package kexec

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"path"
)

func StartCmdListener(files []utils.FilesInfo, zlog *zap.Logger, errChan chan error,matchChan chan *events.KprobeEvent) error {
	go func(errChan chan error) {
		// cleanup old probes
		if err := goebpf.CleanupProbes(); err != nil {
			errChan <- fmt.Errorf("failed to clean probs %s", err.Error())
			return
		}
		ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder()
		if err != nil {
			errChan <- err
			return
		}

		filePath := path.Join(ebpfCompiledFolder, files[0].Name)
		p, err := LoadProgram(filePath)
		if err != nil {
			errChan <- fmt.Errorf("failed to load ebpf program %s", err.Error())
			return
		}
		p.ShowInfo()
		// attach ebpf kprobes
		if err := p.AttachProbes(matchChan); err != nil {
			zlog.Error(fmt.Sprintf("Attach Probes failed: %s", err.Error()))
		}
		defer func() {
			err := p.DetachProbes()
			if err != nil {
				zlog.Error(fmt.Sprintf("Detach Probes failed: %s", err.Error()))
			}
		}()
		ctrlC := make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt)
		<-ctrlC
		// display some stats
		zlog.Info(fmt.Sprintf("%d Event(s) Received\n", p.EventsReceived()))
		zlog.Info(fmt.Sprintf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost()))
		zlog.Info(fmt.Sprintf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost()))
	}(errChan)
	return <-errChan
}
