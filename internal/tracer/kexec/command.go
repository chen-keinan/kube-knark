package kexec

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"go.uber.org/zap"
	"path"
)

//StartCmdListener start exec listener for exec program events
func StartCmdListener(files []utils.FilesInfo, zlog *zap.Logger, errChan chan bool, matchChan chan *events.KprobeEvent) {
	go func(errChan chan bool) {
		// cleanup old probes
		if err := goebpf.CleanupProbes(); err != nil {
			return
		}
		ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder()
		if err != nil {
			return
		}
		filePath := path.Join(ebpfCompiledFolder, files[0].Name)
		p, err := LoadProgram(filePath)
		if err != nil {
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
		<-errChan
	}(errChan)
}
