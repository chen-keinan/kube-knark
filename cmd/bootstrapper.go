package cmd

import (
	"context"
	"fmt"
	"github.com/chen-keinan/go-user-plugins/uplugin"
	"github.com/chen-keinan/kube-knark/internal/common"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/internal/tracer/kexec"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/internal/workers"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
	"github.com/chen-keinan/kube-knark/pkg/ui"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"plugin"
)

// StartKnark start kube-knark event tracer
func StartKnark() {
	app := fx.New(
		// dependency injection
		fx.Provide(NewZapLogger),
		// validation spec files
		fx.Provide(utils.NewKFolder),
		fx.Provide(provideSpecFiles),
		fx.Provide(provideSpecRoutes),
		fx.Provide(provideAPISpecMap),
		fx.Provide(provideFSSpecMap),
		fx.Provide(provideFSSpecCache),
		fx.Provide(mux.NewRouter),
		fx.Provide(matches.NewRouteMatches),
		fx.Provide(utils.GetEbpfCompiledFolder),
		fx.Provide(shell.NewClangCompiler),
		fx.Provide(provideCompiledFiles),
		// init cmd workers
		fx.Provide(LoadAPICallPluginSymbols),
		fx.Provide(LoadFileChangePluginSymbols),
		fx.Provide(ui.NewNetEvtChan),
		fx.Provide(ui.NewKubeKnarkUI),
		fx.Provide(numOfWorkers),
		fx.Provide(matchCmdChan),
		fx.Provide(workers.NewCommandMatchesWorker),
		// init packet workers
		fx.Provide(matchNetChan),
		fx.Provide(workers.NewPacketMatchesWorker),
		fx.Provide(workers.NewPacketMatchData),
		fx.Provide(matches.NewFSMatches),
		fx.Provide(workers.NewCommandMatchesData),
		fx.Provide(ui.NewFilesystemEvtChan),
		fx.Invoke(runKnarkService),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

// load ebpf program and trace events
func runKnarkService(lifecycle fx.Lifecycle,
	netUIChan chan model.K8sAPICallEvent,
	fsUIChan chan model.K8sConfigFileChangeEvent,
	files []utils.FilesInfo,
	NetEventChan chan *netevent.HTTPNetData,
	cmdEventChan chan *execevent.KprobeEvent,
	cm *workers.CommandMatchesWorker,
	pm *workers.PacketMatchesWorker) {

	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		quitChan := make(chan bool)
		errNetChan := make(chan error)
		errCmdChan := make(chan error)
		// invoke cmd msg processing worker
		cm.Invoke()
		// invoke net msg processing worker
		pm.Invoke()
		// start Net Listener
		khttp.StartNetListener(errNetChan, NetEventChan)
		// start exec Listener
		kexec.StartCmdListener(files, errCmdChan, quitChan, cmdEventChan)
		ui.NewKubeKnarkUI(netUIChan, fsUIChan).Draw(errNetChan)
		// wait until Ctrl+C pressed
		ctrlC := make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt)
		select {
		case <-ctrlC:
			// release cmd go routine before panic
			quitChan <- true
			return nil
		case cmdErr := <-errCmdChan:
			panic(cmdErr)
		case netErr := <-errNetChan:
			// release cmd go routine before panic
			quitChan <- true
			panic(netErr)
		}
	},
	})
}

//NewZapLogger zap logger object
func NewZapLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create zap logger instance")
	}
	return logger
}

//matchNetChan return channel for net packet match
func matchNetChan() chan *netevent.HTTPNetData {
	return make(chan *netevent.HTTPNetData, 1000)
}

//matchCmdChan return channel for cmd packet match
func matchCmdChan() chan *execevent.KprobeEvent {
	return make(chan *execevent.KprobeEvent, 1000)
}

//numOfWorkers return num of cmd workers
func numOfWorkers() int {
	return 15
}

//provideCompiledFiles return ebpf compiled files
func provideCompiledFiles(sc shell.ClangExecutor, folder string) []utils.FilesInfo {
	fm := utils.NewKFolder()
	err := utils.CreateKubeKnarkFolders(fm)
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
	fm := utils.NewKFolder()
	err := utils.CreateKubeKnarkFolders(fm)
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
	folder, err := utils.GetSpecAPIFolder(fm)
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
func provideSpecRoutes(files []string) []specs.Routes {
	routesFile, err := specs.BuildSpecRoutes(files)
	if err != nil {
		panic(err)
	}
	return routesFile
}

//provideAPISpecMap provide spec api cache for endpoint validation
func provideAPISpecMap(files []string) map[string]*specs.API {
	specMap, err := specs.CreateMapFromSpecFiles(files)
	if err != nil {
		panic(err)
	}
	return specMap
}

//provideFSSpecMap provide spec fs map validation
func provideFSSpecMap() map[string]interface{} {
	dataFiles, err := getDataFileContent()
	if err != nil {
		panic(err)
	}
	specMap, err := specs.CreateFSMapFromSpecFiles(dataFiles)
	if err != nil {
		panic(err)
	}
	return specMap
}

//provideFSSpecCache provide spec fs cache
func provideFSSpecCache() map[string]*specs.FS {
	dataFiles, err := getDataFileContent()
	if err != nil {
		panic(err)
	}
	specMap, err := specs.CreateFSCacheFromSpecFiles(dataFiles)
	if err != nil {
		panic(err)
	}
	return specMap
}

func getDataFileContent() ([]string, error) {
	fi, err := startup.GenerateFileSystemSpec()
	if err != nil {
		panic(err)
	}
	err = startup.SaveFilesIfNotExist(fi, utils.GetSpecFilesystemFolder)
	if err != nil {
		panic(err)
	}
	folder, err := utils.GetSpecFilesystemFolder(utils.NewKFolder())
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
	return dataFiles, err
}

//LoadAPICallPluginSymbols load API call plugin symbols
func LoadAPICallPluginSymbols(log *zap.Logger) workers.K8sAPICallHook {
	fm := utils.NewKFolder()
	sourceFolder, err := utils.GetPluginSourceSubFolder(fm)
	if err != nil {
		panic(fmt.Sprintf("failed tpo get plugin source sourceFolder %s", err.Error()))
	}
	compliledFolder, err := utils.GetCompilePluginSubFolder(fm)
	if err != nil {
		panic(fmt.Sprintf("failed to get plugin compiled sourceFolder %s", err.Error()))
	}

	pl := uplugin.NewPluginLoader(sourceFolder, compliledFolder)
	names, err := pl.Plugins(uplugin.CompiledExt)
	if err != nil {
		panic(fmt.Sprintf("failed to get plugin compiled plugins %s", err.Error()))
	}
	apiPlugin := workers.K8sAPICallHook{Plugins: make([]plugin.Symbol, 0), Plug: pl}
	for _, name := range names {
		sym, err := pl.Load(name, common.OnK8sAPICallHook)
		if err != nil {
			log.Error(fmt.Sprintf("failed to load sym %s error %s", name, err.Error()))
			continue
		}
		apiPlugin.Plugins = append(apiPlugin.Plugins, sym)
	}
	return apiPlugin
}

//LoadFileChangePluginSymbols load config file change plugin symbols
func LoadFileChangePluginSymbols(log *zap.Logger) workers.K8sFileConfigChangeHook {
	fm := utils.NewKFolder()
	sourceFolder, err := utils.GetPluginSourceSubFolder(fm)
	if err != nil {
		panic(fmt.Sprintf("failed tpo get plugin source sourceFolder %s", err.Error()))
	}
	compliledFolder, err := utils.GetCompilePluginSubFolder(fm)
	if err != nil {
		panic(fmt.Sprintf("failed to get plugin compiled sourceFolder %s", err.Error()))
	}

	pl := uplugin.NewPluginLoader(sourceFolder, compliledFolder)
	names, err := pl.Plugins(uplugin.CompiledExt)
	if err != nil {
		panic(fmt.Sprintf("failed to get plugin compiled plugins %s", err.Error()))
	}
	filePlugin := workers.K8sFileConfigChangeHook{Plugins: make([]plugin.Symbol, 0), Plug: pl}
	for _, name := range names {
		sym, err := pl.Load(name, common.OnK8sFileConfigChangeHook)
		if err != nil {
			log.Error(fmt.Sprintf("failed to load sym %s error %s", name, err.Error()))
			continue
		}
		filePlugin.Plugins = append(filePlugin.Plugins, sym)
	}
	return filePlugin
}
