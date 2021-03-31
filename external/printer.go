//nolint
package external

import (
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"io"
	"os"
	"sync"
)

var WaitGroup sync.WaitGroup
var PrinterWaitGroup sync.WaitGroup

// Printer output parsed http messages
type Printer struct {
	OutputQueue chan *netevent.HTTPNetData
	outputFile  io.WriteCloser
}

var maxOutputQueueLen = 4096

func NewPrinter(matchChan chan *netevent.HTTPNetData) *Printer {
	var outputFile io.WriteCloser
	outputFile = os.Stdout
	printer := &Printer{OutputQueue: make(chan *netevent.HTTPNetData, maxOutputQueueLen), outputFile: outputFile}
	printer.start(matchChan)
	return printer
}

func (p *Printer) Send(data *netevent.HTTPNetData) {
	if len(data.HTTPRequestData.Method) == 0 ||
		len(data.HTTPRequestData.Host) == 0 ||
		len(data.HTTPRequestData.RequestURI) == 0 {
		return
	}
	p.OutputQueue <- data
}

func (p *Printer) start(matchChan chan *netevent.HTTPNetData) {
	PrinterWaitGroup.Add(1)
	go p.printBackground(matchChan)
}

func (p *Printer) printBackground(matchChan chan *netevent.HTTPNetData) {
	defer PrinterWaitGroup.Done()
	defer p.outputFile.Close()
	for msg := range p.OutputQueue {
		matchChan <- msg
	}
}

func (p *Printer) Finish() {
	close(p.OutputQueue)
}
