package kexec

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"path"
)

//StartCmdListener start exec listener for exec program events
func StartCmdListener(files []utils.FilesInfo, errChan chan error, quitChan chan bool, cmdEventChan chan *events.KprobeEvent) {
	go func(quitChan chan bool, errChan chan error) {
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
		if err := p.AttachProbes(cmdEventChan); err != nil {
			errChan <- fmt.Errorf("Attach Probes failed: %s", err.Error())
		}
		defer func() {
			err := p.DetachProbes()
			if err != nil {
				errChan <- fmt.Errorf("Detach Probes failed: %s", err.Error())
			}
		}()
		<-quitChan
	}(quitChan, errChan)
}
