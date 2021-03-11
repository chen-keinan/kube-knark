package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/trace"
	"os"
	"os/signal"
)

func StartKnark() {
	//	log := logger.NewZapLogger()
	flag.Parse()

	// Start the exec syscall monitor.
	m, err := trace.NewMonitor()
	if err != nil {

	}

	done := make(chan struct{})
	events, err := m.Start(done)
	if err != nil {
	}

	// Handle signals for shutting down.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		close(done)
	}()

	// Read incoming exec events.
	for e := range events {
		data, _ := json.Marshal(e)
		fmt.Println(string(data))
	}
}
