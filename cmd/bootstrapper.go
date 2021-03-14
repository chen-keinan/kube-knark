package cmd

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/trace"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"log"
	"os"
	"os/signal"
	"path"
)

func StartKnark() {
	/*fileInfo, err := startup.GenerateEbpfFiles()
	if err != nil {
		panic(err)
	}
	err = startup.SaveEbpfFilesIfNotExist(fileInfo)
	if err != nil {
		panic(err)
	}
	err = startup.CompileEbpfSources(fileInfo)
	if err != nil {
		panic(err)
	}
	// cleanup old probes
	if err := goebpf.CleanupProbes(); err != nil {
		log.Println(err)
	}*/

	// load ebpf program

	if err := goebpf.CleanupProbes(); err != nil {
		log.Println(err)
	}
	ebpfCompiledFolder := utils.GetEbpfCompiledFolder()
	files, err := utils.GetEbpfFiles(ebpfCompiledFolder)
	filePath := path.Join(ebpfCompiledFolder, files[0].Name)
	fmt.Print(filePath)
	p, err := trace.LoadProgram(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to load ebpf program %s", err.Error()))
	}
	p.ShowInfo()

	// attach ebpf kprobes
	if err := p.AttachProbes(); err != nil {
		log.Fatalf("AttachProbes() failed: %v", err)
	}
	defer p.DetachProbes()
	// wait until Ctrl+C pressed
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	<-ctrlC
	// display some stats
	fmt.Println()
	fmt.Printf("%d Event(s) Received\n", p)
	fmt.Printf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost())
	fmt.Printf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost())

}
