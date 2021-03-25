package ui

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	tb "github.com/nsf/termbox-go"
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
		// draw net event and fs event sections
		_, netSection := drawSections(termWidth, termHeight)
		///
		fsTable := widgets.NewTable()
		fsEvents := make([][]string, 0)
		fsEvents = append(fsEvents, []string{"Severity", "Name", "Command args"})
		fsTable.Rows = fsEvents
		fsTable.TextStyle = ui.NewStyle(ui.ColorWhite)
		fsTable.SetRect(1, 1, termWidth-1, termHeight/2)
		fsTable.Title = "K8s configuration file change events"
		ui.Render(fsTable)
		///q
		// render to ui
		ui.Render(p, fsTable, netSection)
		uiEvents := ui.PollEvents()

		netEvents := make([]string, 0)
		for {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					tb.Close()
					return
				}
			case fsEvent := <-kku.FsEvtChan:
				args := fmt.Sprintf("%s", fsEvent.Msg.Args)
				fsEvents = append(fsEvents, []string{fsEvent.Spec.Severity, fsEvent.Spec.Name, args})
				fsTable.Rows = fsEvents
				ui.Render(fsTable)

			case netEvent := <-kku.NetEvtChan:
				value := fmt.Sprintf("%s:%s:%s %s", netEvent.Spec.Severity,
					netEvent.Spec.Name,
					netEvent.Msg.HTTPRequestData.Method,
					netEvent.Msg.HTTPRequestData.RequestURI)
				netEvents = append(netEvents, value)
				netSection.Rows = netEvents
				ui.Render(netSection)
			}
		}
	}()
}

func drawParagraph(termWidth, termHeight int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "Kube-Knark Tracer"
	p.SetRect(0, 0, termWidth, termHeight)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}

func drawSections(w, h int) (*widgets.List, *widgets.List) {
	//a := createSectionList(1, 1, w-1, h/2, "K8s configuration file change events")
	b := createSectionList(1, h/2, w-1, h-1, "K8s API change events")
	return nil, b
}

func createSectionList(x0, y0, x1, y1 int, title string) *widgets.List {
	l := widgets.NewList()
	l.Title = title
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = true
	l.SetRect(x0, y0, x1, y1)
	return l
}
