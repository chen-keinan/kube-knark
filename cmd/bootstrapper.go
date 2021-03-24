package cmd

import (
	"context"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/logger"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/internal/tracer/kexec"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/internal/workers"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

// StartKnark start kube-knark event tracer
func StartKnark() {
	app := fx.New(
		// dependency injection
		fx.Provide(logger.NewZapLogger),
		// validation spec files
		fx.Provide(provideSpecFiles),
		fx.Provide(provideSpecRoutes),
		fx.Provide(provideAPISpecMap),
		fx.Provide(provideFSSpecMap),
		fx.Provide(mux.NewRouter),
		fx.Provide(matches.NewRouteMatches),
		fx.Provide(utils.GetEbpfCompiledFolder),
		fx.Provide(shell.NewClangCompiler),
		fx.Provide(provideCompiledFiles),
		// init cmd workers
		fx.Provide(numOfWorkers),
		fx.Provide(matchCmdChan),
		fx.Provide(workers.NewCommandMatchesWorker),
		// init packet workers
		fx.Provide(matchNetChan),
		fx.Provide(workers.NewPacketMatchesWorker),
		fx.Provide(providePacketData),
		fx.Provide(provideFSMatches),
		fx.Provide(provideCommandData),
		fx.Invoke(runKnarkService),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

// load ebpf program and trace events
func runKnarkService(lifecycle fx.Lifecycle,
	zlog *zap.Logger,
	files []utils.FilesInfo,
	NetEventChan chan *khttp.HTTPNetData,
	cmdEventChan chan *events.KprobeEvent,
	cm *workers.CommandMatchesWorker,
	pm *workers.PacketMatchesWorker) {

	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		quitChan := make(chan bool)
		errNetChan := make(chan error)
		errCmdChan := make(chan error)
		cm.Invoke()
		pm.Invoke()
		// start Net Listener
		khttp.StartNetListener(errNetChan, NetEventChan)
		// start exec Listener
		kexec.StartCmdListener(files, errCmdChan, quitChan, cmdEventChan)
		// wait until Ctrl+C pressed
		ctrlC := make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt)
		select {
		case <-ctrlC:
			return nil
		case cmdErr := <-errCmdChan:
			panic(cmdErr)
		case netErr := <-errNetChan:
			panic(netErr)
		}
	},
	})
}

//matchNetChan return channel for net packet match
func matchNetChan() chan *khttp.HTTPNetData {
	return make(chan *khttp.HTTPNetData, 1000)
}

//matchCmdChan return channel for cmd packet match
func matchCmdChan() chan *events.KprobeEvent {
	return make(chan *events.KprobeEvent, 1000)
}

//numOfWorkers return num of cmd workers
func numOfWorkers() int {
	return 15
}

//provideCompiledFiles return ebpf compiled files
func provideCompiledFiles(sc shell.ClangExecutor, folder string) []utils.FilesInfo {
	err := utils.CreateKubeKnarkFolders()
	if err != nil {
		panic(err)
	}
	fi, err := startup.GenerateEbpfFiles()
	if err != nil {
		panic(err)
	}
	err = startup.SaveFilesIfNotExist(fi, utils.GetEbpfSourceFolder)
	if err != nil {
		panic(err)
	}
	err = startup.CompileEbpfSources(fi, sc)
	if err != nil {
		panic(err)
	}
	files, err := utils.GetFiles(folder)
	if err != nil {
		panic(err)
	}
	return files
}

//provideSpecFiles return spec files
func provideSpecFiles() []string {
	err := utils.CreateKubeKnarkFolders()
	if err != nil {
		panic(err)
	}
	fi, err := startup.GenerateSpecFiles()
	if err != nil {
		panic(err)
	}
	err = startup.SaveFilesIfNotExist(fi, utils.GetSpecAPIFolder)
	if err != nil {
		panic(err)
	}
	folder, err := utils.GetSpecAPIFolder()
	if err != nil {
		panic(err)
	}
	files, err := utils.GetFiles(folder)
	if err != nil {
		panic(err)
	}
	dataFiles := make([]string, 0)
	for _, f := range files {
		dataFiles = append(dataFiles, f.Data)
	}
	return dataFiles
}

//provideSpecRoutes provide spec api route for endpoint validation
func provideSpecRoutes(files []string) []routes.Routes {
	routesFile, err := routes.BuildSpecRoutes(files)
	if err != nil {
		panic(err)
	}
	return routesFile
}

//provideAPISpecMap provide spec api cache for endpoint validation
func provideAPISpecMap(files []string) map[string]*routes.API {
	specMap, err := routes.CreateMapFromSpecFiles(files)
	if err != nil {
		panic(err)
	}
	return specMap
}

//provideFSSpecMap provide spec fs map validation
func provideFSSpecMap() map[string]interface{} {
	fi, err := startup.GenerateFileSystemSpec()
	if err != nil {
		panic(err)
	}
	err = startup.SaveFilesIfNotExist(fi, utils.GetSpecFilesystemFolder)
	if err != nil {
		panic(err)
	}
	folder, err := utils.GetSpecFilesystemFolder()
	if err != nil {
		panic(err)
	}
	files, err := utils.GetFiles(folder)
	if err != nil {
		panic(err)
	}
	dataFiles := make([]string, 0)
	for _, f := range files {
		dataFiles = append(dataFiles, f.Data)
	}
	specMap, err := routes.CreateFSMapFromSpecFiles(dataFiles)
	if err != nil {
		panic(err)
	}
	return specMap
}

//providePacketData provide spec data for packet worker
func providePacketData(rm *matches.RouteMatches, pmc chan *khttp.HTTPNetData, cache map[string]*routes.API, numOfWorkers int) *workers.PacketMatchData {
	return workers.NewPacketMatchData(rm, pmc, cache, numOfWorkers)
}

//provideFSMatches return fs matches instance
func provideFSMatches(fsCommandMap map[string]interface{}) *matches.FSMatches {
	return matches.NewFSMatches(fsCommandMap)
}

//provideCommandData provide spec data for command worker
func provideCommandData(cmc chan *events.KprobeEvent, NumOfWorkers int, fsMatches *matches.FSMatches) *workers.CommandMatchData {
	return workers.NewCommandMatchesData(cmc, NumOfWorkers, fsMatches)
}
