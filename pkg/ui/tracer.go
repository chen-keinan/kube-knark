package ui

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

//KubeKnarkUI return UI object
type KubeKnarkUI struct {
	NetEvtChan chan NetEvt
	FsEvtChan  chan FilesystemEvt
}

// NetEvt net event msg
type NetEvt struct {
	Msg  *khttp.HTTPNetData
	Spec *routes.API
}

// FilesystemEvt fs event msg
type FilesystemEvt struct {
	Msg  *events.KprobeEvent
	Spec *routes.FS
}

//NewNetEvtChan return net event channel
func NewNetEvtChan() chan NetEvt {
	return make(chan NetEvt)
}

//NewFilesystemEvtChan return file system event channel
func NewFilesystemEvtChan() chan FilesystemEvt {
	return make(chan FilesystemEvt)
}

//NewKubeKnarkUI return new KubeKnarkUI object
func NewKubeKnarkUI(netData chan NetEvt, fsData chan FilesystemEvt) *KubeKnarkUI {
	return &KubeKnarkUI{NetEvtChan: netData, FsEvtChan: fsData}
}

//Draw draw ui kube knark ui with paragraph and lists
func (kku *KubeKnarkUI) Draw(errNetChan chan error) {
	go func() {
		if err := ui.Init(); err != nil {
			errNetChan <- err
			return
		}
		defer ui.Close()
		// draw external paragraph
		termWidth, termHeight := ui.TerminalDimensions()
		p := drawParagraph(termWidth, termHeight)
		// init event tables
		fsTable, fsEvents := kku.drawFileSystemTable(termWidth, termHeight)
		netTable, netEvents := kku.drawNetTable(termWidth, termHeight)
		// render to ui
		ui.Render(p, fsTable, netTable)
		uiEvents := ui.PollEvents()
		// watch for net , file system events
		kku.watchEvents(uiEvents, fsTable, netTable, fsEvents, netEvents)
	}()
}

func (kku *KubeKnarkUI) watchEvents(uiEvents <-chan ui.Event, fsTable *Table, netTable *Table, fsEvents [][]string, netEvents [][]string) {
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "j", "<Down>":
				fsTable.ScrollDown()
			case "k", "<Up>":
				fsTable.ScrollUp()
			case "w":
				netTable.ScrollUp()
			case "s":
				fsTable.ScrollDown()
			}
		case fsEvent := <-kku.FsEvtChan:
			args := fmt.Sprintf("%s", fsEvent.Msg.Args)
			fsEvents = append(fsEvents, []string{fsEvent.Spec.Severity, fsEvent.Spec.Name, args})
			fsTable.Rows = fsEvents
			ui.Render(fsTable)

		case netEvent := <-kku.NetEvtChan:
			netEvents = append(netEvents, []string{netEvent.Spec.Severity, netEvent.Spec.Name, netEvent.Msg.HTTPRequestData.RequestURI})
			netTable.Rows = netEvents
			ui.Render(netTable)
		}
		ui.Render(fsTable)
		ui.Render(netTable)
	}
}

func (kku *KubeKnarkUI) drawFileSystemTable(termWidth int, termHeight int) (*Table, [][]string) {
	fsTable := NewTable(true)
	fsEvents := make([][]string, 0)
	fsEvents = append(fsEvents, []string{"Severity", "Name", "Command args"})
	fsTable.Rows = fsEvents
	fsTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	fsTable.SetRect(1, 1, termWidth-1, termHeight/2)
	fsTable.Title = "K8s configuration file change events"
	ui.Render(fsTable)
	return fsTable, fsEvents
}

func (kku *KubeKnarkUI) drawNetTable(termWidth int, termHeight int) (*Table, [][]string) {
	netTable := NewTable(true)
	netEvents := make([][]string, 0)
	netEvents = append(netEvents, []string{"Severity", "Name", "API Call"})
	netTable.Rows = netEvents
	netTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	netTable.SetRect(1, termHeight/2, termWidth-1, termHeight-1)
	netTable.Title = "K8s API change events"
	ui.Render(netTable)
	return netTable, netEvents
}

func drawParagraph(termWidth, termHeight int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "Kube-Knark Tracer"
	p.SetRect(0, 0, termWidth, termHeight)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}
