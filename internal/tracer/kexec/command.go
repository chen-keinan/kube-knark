package kexec

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"path"
)

//StartCmdListener start exec listener for exec program events
func StartCmdListener(files []utils.FilesInfo, errChan chan error, quitChan chan bool, cmdEventChan chan *execevent.KprobeEvent) {
	go func(quitChan chan bool, errChan chan error) {
		// cleanup old probes
		if err := goebpf.CleanupProbes(); err != nil {
			errChan <- fmt.Errorf("cleanup Probes failed: %s", err.Error())
		}
		ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder(utils.NewKFolder())
		if err != nil {
			errChan <- fmt.Errorf("get Ebpf program failed: %s", err.Error())
		}
		filePath := path.Join(ebpfCompiledFolder, files[0].Name)
		p, err := LoadProgram(filePath)
		if err != nil {
			errChan <- fmt.Errorf("load ebpf program failed: %s", err.Error())
		}
		p.ShowInfo()
		// attach ebpf kprobes
		if err := p.AttachProbes(cmdEventChan); err != nil {
			errChan <- fmt.Errorf("attach probes failed: %s", err.Error())
		}
		defer func() {
			err := p.DetachProbes()
			if err != nil {
				errChan <- fmt.Errorf("detach probes failed: %s", err.Error())
			}
		}()
		<-quitChan
	}(quitChan, errChan)
}
