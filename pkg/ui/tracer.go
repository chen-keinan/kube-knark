package ui

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"sort"
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
func (kui *KubeKnarkUI) Draw(errNetChan chan error) {
	go func() {
		if err := ui.Init(); err != nil {
			errNetChan <- err
			return
		}
		defer ui.Close()
		// draw external paragraph
		termWidth, termHeight := ui.TerminalDimensions()
		p := kui.buildParagraph(termWidth, termHeight)
		// init event tables
		fsTable, headers := kui.buildFileSystemTable(termWidth, termHeight)
		netTable, netEvents := kui.buildNetTable(termWidth, termHeight)
		// render to ui
		ui.Render(p, fsTable, netTable)
		uiEvents := ui.PollEvents()
		// watch for net , file system events
		kui.watchEvents(uiEvents, fsTable, netTable, headers, netEvents)
	}()
}

func (kui *KubeKnarkUI) watchEvents(uiEvents <-chan ui.Event, fsTable *Table, netTable *Table, fsHeaders []string, netHeaders []string) {
	fsEvts := make([]*FilesystemEvt, 0)
	netEvts := make([]*NetEvt, 0)
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
				netTable.ScrollDown()
			}
		case fsEvent := <-kui.FsEvtChan:
			var fsRows = [][]string{fsHeaders}
			fsTable.Rows = kui.sortFSRows(&fsEvts, &fsEvent, fsRows)

		case netEvent := <-kui.NetEvtChan:
			var netRow = [][]string{netHeaders}
			netTable.Rows = kui.sortNetRows(&netEvts, &netEvent, netRow)
		}
		ui.Render(fsTable)
		ui.Render(netTable)
	}
}

func (kui *KubeKnarkUI) sortNetRows(netEvts *[]*NetEvt, netEvent *NetEvt, netRows [][]string) [][]string {
	*netEvts = append(*netEvts, netEvent)
	// sort event by severity
	sort.Slice(*netEvts, func(i, j int) bool {
		return (*netEvts)[i].Spec.SeverityInt < (*netEvts)[j].Spec.SeverityInt
	})
	for _, nse := range *netEvts {
		netRows = append(netRows, []string{nse.Spec.Severity, nse.Spec.Name, nse.Msg.HTTPRequestData.Method, nse.Msg.HTTPRequestData.RequestURI, nse.Msg.HTTPRequestData.StartTime})
	}
	return netRows
}

// sort table by severity
func (kui *KubeKnarkUI) sortFSRows(fsEvts *[]*FilesystemEvt, fsEvent *FilesystemEvt, fsRows [][]string) [][]string {
	*fsEvts = append(*fsEvts, fsEvent)
	// sort event by severity
	sort.Slice(*fsEvts, func(i, j int) bool {
		return (*fsEvts)[i].Spec.SeverityInt < (*fsEvts)[j].Spec.SeverityInt
	})
	for _, fse := range *fsEvts {
		args := fmt.Sprintf("%s", fse.Msg.Args)
		fsRows = append(fsRows, []string{fse.Spec.Severity, fse.Spec.Name, args, fse.Msg.StartTime})
	}
	return fsRows
}

func (kui *KubeKnarkUI) buildFileSystemTable(termWidth int, termHeight int) (*Table, []string) {
	fsTable := NewTable(true)
	longColumn := (termWidth - 45) / 2
	fsTable.Table.ColumnWidths = []int{10, longColumn - 15, longColumn + 7, 40}
	fsRows := make([][]string, 0)
	headers := []string{"Severity", "Name", "Command args", "Created"}
	fsRows = append(fsRows, headers)
	fsTable.Rows = fsRows
	fsTable.RowSeparator = false
	fsTable.Colors.SelectedRowBg = ui.ColorGreen
	fsTable.Colors.Text = ui.ColorWhite
	fsTable.Colors.SelectedRowBg = ui.ColorBlue
	fsTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	fsTable.SetRect(1, 1, termWidth-1, termHeight/2)
	fsTable.Title = "K8s configuration file change events"
	return fsTable, headers
}

func (kui *KubeKnarkUI) buildNetTable(termWidth int, termHeight int) (*Table, []string) {
	netTable := NewTable(true)
	longColumn := (termWidth - 45) / 2
	netTable.Table.ColumnWidths = []int{10, longColumn - 15, 7, longColumn, 40}
	netRows := make([][]string, 0)
	headers := []string{"Severity", "Name", "Method", "API Call", "Created"}
	netRows = append(netRows, headers)
	netTable.Rows = netRows
	netTable.RowSeparator = false
	netTable.Colors.SelectedRowBg = ui.ColorGreen
	netTable.Colors.Text = ui.ColorWhite
	netTable.Colors.SelectedRowBg = ui.ColorBlue
	netTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	netTable.SetRect(1, termHeight/2, termWidth-1, termHeight-1)
	netTable.Title = "K8s API change events"
	return netTable, headers
}

// draw paragraph section
func (kui *KubeKnarkUI) buildParagraph(termWidth, termHeight int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "Kube-Knark Tracer"
	p.SetRect(0, 0, termWidth, termHeight)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}
