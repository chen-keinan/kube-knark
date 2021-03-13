package cmd

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/trace"
	"github.com/dropbox/goebpf"
	"log"
	"os"
	"os/signal"
)

func StartKnark() {
	// cleanup old probes
	if err := goebpf.CleanupProbes(); err != nil {
		log.Println(err)
	}

	// load ebpf program
	p, err := trace.LoadProgram(common.KprobeFile)
	if err != nil {
		panic("failed to load ebpf program")
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
