package cmd

import (
	"context"
	"fmt"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/ebpf_mgr"
	"github.com/chen-keinan/kube-knark/internal/logger"
	"github.com/chen-keinan/kube-knark/internal/startup"
 	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"path"
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
		// cleanup old probes
		if err := goebpf.CleanupProbes(); err != nil {
			return fmt.Errorf("failed to clean probs %s", err.Error())
		}
		ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder()
		if err != nil {
			return err
		}

		filePath := path.Join(ebpfCompiledFolder, files[0].Name)
		p, err := ebpf_mgr.LoadProgram(filePath)
		if err != nil {
			return fmt.Errorf("failed to load ebpf program %s", err.Error())
		}
		p.ShowInfo()

		// attach ebpf kprobes
		if err := p.AttachProbes(); err != nil {
			zlog.Error(fmt.Sprintf("Attach Probes failed: %s", err.Error()))
		}
		defer func() {
			err := p.DetachProbes()
			if err != nil {
				zlog.Error(fmt.Sprintf("Detach Probes failed: %s", err.Error()))
			}
		}()

		// wait until Ctrl+C pressed
		ctrlC := make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt)
		<-ctrlC

		// display some stats
		zlog.Info(fmt.Sprintf("%d Event(s) Received\n", p.EventsReceived()))
		zlog.Info(fmt.Sprintf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost()))
		zlog.Info(fmt.Sprintf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost()))
		return nil
	},
	})
}
