package cmd

import (
	"context"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/logger"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/internal/tracer/kexec"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/internal/workers"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

// StartKnark start kube-knark event tracer
func StartKnark() {
	app := fx.New(
		fx.Provide(logger.NewZapLogger),
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
	NetChan chan string,
	cmdChan chan *events.KprobeEvent,
	cm *workers.CommandMatches,
	pm *workers.PacketMatches) {

	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		errChan := make(chan error)
		cm.Invoke()
		pm.Invoke()
		// start Net Listener
		go func() {
			khttp.StartNetListener(zlog, NetChan)
		}()
		// start exec Listener
		err := kexec.StartCmdListener(files, zlog, errChan, cmdChan)
		if err != nil {
			zlog.Error(err.Error())
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
func MatchNetChan() chan string {
	return make(chan string, 1000)
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
func ProvideCompiledFiles(sc *shell.ClangCompiler, folder string) []utils.FilesInfo {
	fi, err := startup.GenerateEbpfFiles()
	if err != nil {
		panic(err)
	}
	err = startup.SaveEbpfFilesIfNotExist(fi)
	if err != nil {
		panic(err)
	}
	err = startup.CompileEbpfSources(fi, sc)
	if err != nil {
		panic(err)
	}
	files, err := utils.GetEbpfFiles(folder)
	if err != nil {
		panic(err)
	}
	return files
}
