package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	tb "github.com/nsf/termbox-go"
)

//KubeKnarkUI return UI object
type KubeKnarkUI struct {
	netEvtChan chan NetEvt
	fsEvtChan  chan FilesystemEvt
}

// NetEvt net event msg
type NetEvt struct {
	Msg string
}

// FilesystemEvt fs event msg
type FilesystemEvt struct {
	Msg string
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
	return &KubeKnarkUI{netEvtChan: netData, fsEvtChan: fsData}
}

//Draw draw ui kube knark ui with paragraph and lists
func (kku *KubeKnarkUI) Draw(errNetChan chan error) {
	go func() {
		if err := ui.Init(); err != nil {
			errNetChan <- err
			return
		}
		defer ui.Close()
		// init term buffer
		if err := tb.Init(); err != nil {
			errNetChan <- err
			return
		}
		if err := tb.Sync(); err != nil {
			errNetChan <- err
			return
		}
		// draw external paragraph
		p := drawParagraph()
		// draw net event and fs event sections
		fsSection, netSection := drawSections()
		// render to ui
		ui.Render(p, fsSection, netSection)
		uiEvents := ui.PollEvents()
		fsEvents := make([]string, 0)
		netEvents := make([]string, 0)
		for {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					tb.Close()
					return
				}
			case fsEvent := <-kku.fsEvtChan:
				fsEvents = append(fsEvents, fsEvent.Msg)
				fsSection.Rows = fsEvents
				ui.Render(fsSection)

			case netEvent := <-kku.netEvtChan:
				netEvents = append(netEvents, netEvent.Msg)
				netSection.Rows = netEvents
				ui.Render(netSection)
			}
		}
	}()
}

func drawParagraph() *widgets.Paragraph {
	w, h := tb.Size()
	p := widgets.NewParagraph()
	p.Title = "Kube-Knark Tracer"
	p.SetRect(0, 0, w, h)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}

func drawSections() (*widgets.List, *widgets.List) {
	w, h := tb.Size()
	a := createSectionList(1, 1, w-1, h/2, "K8s configuration file change events")
	b := createSectionList(1, h/2, w-1, h-1, "K8s API change events")
	return a, b
}

func createSectionList(x0, y0, x1, y1 int, title string) *widgets.List {
	l := widgets.NewList()
	l.Title = title
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = true
	l.SetRect(x0, y0, x1, y1)
	return l
}
