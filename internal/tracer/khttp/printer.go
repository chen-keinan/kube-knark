package khttp

import (
	"fmt"
	"io"
	"os"
)

// Printer output parsed http messages
type Printer struct {
	OutputQueue chan string
	outputFile  io.WriteCloser
}

var maxOutputQueueLen = 4096

func newPrinter(matchChan chan string) *Printer {
	var outputFile io.WriteCloser
	outputFile = os.Stdout
	printer := &Printer{OutputQueue: make(chan string, maxOutputQueueLen), outputFile: outputFile}
	printer.start(matchChan)
	return printer
}

func (p *Printer) send(msg string) {
	if len(p.OutputQueue) == maxOutputQueueLen {
		// skip this msg
		fmt.Fprintln(os.Stderr, "too many messages to output, discard current!")
		return
	}
	p.OutputQueue <- msg
}

func (p *Printer) start(matchChan chan string) {
	printerWaitGroup.Add(1)
	go p.printBackground(matchChan)
}

func (p *Printer) printBackground(matchChan chan string) {
	defer printerWaitGroup.Done()
	defer p.outputFile.Close()
	for msg := range p.OutputQueue {
		matchChan <- msg
	}
}

func (p *Printer) finish() {
	close(p.OutputQueue)
}
