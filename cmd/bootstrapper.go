package cmd

import (
	"context"
	"fmt"
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
		fx.Provide(logger.NewZapLogger),
		fx.Provide(ProvideSpecFiles),
		fx.Provide(ProvideSpecRoutes),
		fx.Provide(mux.NewRouter),
		fx.Provide(matches.NewRouteMatches),
		fx.Provide(utils.GetEbpfCompiledFolder),
		fx.Provide(shell.NewClangCompiler),
		fx.Provide(ProvideCompiledFiles),
		// init cmd workers
		fx.Provide(NumOfWorkers),
		fx.Provide(MatchCmdChan),
		fx.Provide(workers.NewCommandMatches),
		// init packet workers
		fx.Provide(MatchNetChan),
		fx.Provide(workers.NewPacketMatches),
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
	NetChan chan *khttp.HTTPNetData,
	cmdChan chan *events.KprobeEvent,
	cm *workers.CommandMatches,
	pm *workers.PacketMatches) {

	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		errChan := make(chan error)
		cm.Invoke()
		pm.Invoke()
		// start Net Listener
		go func() {
			err := khttp.StartNetListener(zlog, NetChan)
			if err != nil {
				panic("failed to init net listener")
			}
		}()
		// start exec Listener
		err := kexec.StartCmdListener(files, zlog, errChan, cmdChan)
		if err != nil {
			panic("failed to init cmd listener")
		}
		// wait until Ctrl+C pressed
		ctrlC := make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt)
		<-ctrlC
		return nil
	},
	})

}

//MatchNetChan return channel for net packet match
func MatchNetChan() chan *khttp.HTTPNetData {
	return make(chan *khttp.HTTPNetData, 1000)
}

//MatchCmdChan return channel for cmd packet match
func MatchCmdChan() chan *events.KprobeEvent {
	return make(chan *events.KprobeEvent, 1000)
}

//NumOfWorkers return num of cmd workers
func NumOfWorkers() int {
	return 5
}

//ProvideCompiledFiles return ebpf compiled files
func ProvideCompiledFiles(sc shell.ClangExecutor, folder string) []utils.FilesInfo {
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

//ProvideSpecFiles return spec files
func ProvideSpecFiles() []string {
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
	fmt.Println(files)
	if err != nil {
		panic(err)
	}
	dataFiles := make([]string, 0)
	for _, f := range files {
		dataFiles = append(dataFiles, f.Data)
	}
	return dataFiles
}

//ProvideSpecRoutes provide spec api route for endpoint validation
func ProvideSpecRoutes(files []string) []routes.Routes {
	routes, err := routes.BuildSpecRoutes(files)
	if err != nil {
		panic(err)
	}
	return routes
}
