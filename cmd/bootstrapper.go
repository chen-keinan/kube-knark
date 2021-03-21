package cmd

import (
	"context"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/logger"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/internal/tracer/kexec"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
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
		fx.Invoke(runKnarkService),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
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

// load ebpf program and trace events
func runKnarkService(lifecycle fx.Lifecycle, zlog *zap.Logger, files []utils.FilesInfo) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		errChan := make(chan error)
		// start Net Listener
		err := khttp.StartNetListener(zlog)
		if err != nil {
			return err
		}
		// start exec Listener
		err = kexec.StartCmdListener(files, zlog, errChan)
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
