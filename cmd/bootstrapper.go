package cmd

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/internal/trace"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/dropbox/goebpf"
	"log"
	"os"
	"os/signal"
)

func StartKnark() {
	fileInfo, err := startup.GenerateEbpfFiles()
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
	}

	// load ebpf program
	files, err := utils.GetEbpfFiles(utils.GetEbpfCompiledFolder())
	if err := goebpf.CleanupProbes(); err != nil {
		log.Println(err)
	}
	for _, ebpfFile := range files {
		go func() {
			p, err := trace.LoadProgram(ebpfFile.Name)
			if err != nil {
				panic("failed to load ebpf program")
			}
			p.ShowInfo()

			// attach ebpf kprobes
			if err := p.AttachProbes(); err != nil {
				log.Fatalf("AttachProbes() failed: %v", err)
			}
			defer p.DetachProbes()

			// display some stats
			fmt.Println()
			fmt.Printf("%d Event(s) Received\n", p)
			fmt.Printf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost())
			fmt.Printf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.EventsLost())
		}()
	}
	// wait until Ctrl+C pressed
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	<-ctrlC
}
