package ui

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"sort"
)

//KubeKnarkUI return UI object
type KubeKnarkUI struct {
	NetEvtChan chan model.K8sAPICallEvent
	FsEvtChan  chan model.K8sConfigFileChangeEvent
	fsTable    *Table
	netTable   *Table
	fsHeaders  []string
	netHeaders []string
	paragraph  *widgets.Paragraph
}

//NewNetEvtChan return net event channel
func NewNetEvtChan() chan model.K8sAPICallEvent {
	return make(chan model.K8sAPICallEvent)
}

//NewFilesystemEvtChan return file system event channel
func NewFilesystemEvtChan() chan model.K8sConfigFileChangeEvent {
	return make(chan model.K8sConfigFileChangeEvent)
}

//NewKubeKnarkUI return new KubeKnarkUI object
func NewKubeKnarkUI(netData chan model.K8sAPICallEvent, fsData chan model.K8sConfigFileChangeEvent) *KubeKnarkUI {
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
		kui.buildParagraph(termWidth, termHeight)
		// init event tables
		kui.buildFileSystemTable(termWidth, termHeight)
		kui.buildNetTable(termWidth, termHeight)
		// render to ui
		ui.Render(kui.paragraph, kui.fsTable, kui.netTable)
		// watch for net , file system events
		kui.watchEvents(ui.PollEvents())
	}()
}

func (kui *KubeKnarkUI) watchEvents(uiEvents <-chan ui.Event) {
	fsEvts := make([]*model.K8sConfigFileChangeEvent, 0)
	netEvts := make([]*model.K8sAPICallEvent, 0)
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "j", "<Down>":
				kui.fsTable.ScrollDown()
			case "k", "<Up>":
				kui.fsTable.ScrollUp()
			case "w":
				kui.netTable.ScrollUp()
			case "s":
				kui.netTable.ScrollDown()
			}
		case fsEvent := <-kui.FsEvtChan:
			var fsRows = [][]string{kui.fsHeaders}
			kui.fsTable.Rows = kui.sortFSRows(&fsEvts, &fsEvent, fsRows)

		case netEvent := <-kui.NetEvtChan:
			var netRow = [][]string{kui.netHeaders}
			kui.netTable.Rows = kui.sortNetRows(&netEvts, &netEvent, netRow)
		}
		ui.Render(kui.fsTable)
		ui.Render(kui.netTable)
	}
}

func (kui *KubeKnarkUI) sortNetRows(netEvts *[]*model.K8sAPICallEvent, netEvent *model.K8sAPICallEvent, netRows [][]string) [][]string {
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
func (kui *KubeKnarkUI) sortFSRows(fsEvts *[]*model.K8sConfigFileChangeEvent, fsEvent *model.K8sConfigFileChangeEvent, fsRows [][]string) [][]string {
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

func (kui *KubeKnarkUI) buildFileSystemTable(termWidth int, termHeight int) {
	kui.fsTable = NewTable(true)
	longColumn := (termWidth - 45) / 2
	kui.fsTable.Table.ColumnWidths = []int{10, longColumn - 15, longColumn + 7, 40}
	fsRows := make([][]string, 0)
	kui.fsHeaders = []string{"Severity", "Name", "command args", "Created"}
	fsRows = append(fsRows, kui.fsHeaders)
	kui.fsTable.Rows = fsRows
	kui.fsTable.RowSeparator = false
	kui.fsTable.Colors.SelectedRowBg = ui.ColorGreen
	kui.fsTable.Colors.Text = ui.ColorWhite
	kui.fsTable.Colors.SelectedRowBg = ui.ColorBlue
	kui.fsTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	kui.fsTable.SetRect(1, 1, termWidth-1, termHeight/2)
	kui.fsTable.Title = "K8s configuration file change events"
}

func (kui *KubeKnarkUI) buildNetTable(termWidth int, termHeight int) {
	kui.netTable = NewTable(true)
	longColumn := (termWidth - 45) / 2
	kui.netTable.Table.ColumnWidths = []int{10, longColumn - 15, 7, longColumn, 40}
	netRows := make([][]string, 0)
	kui.netHeaders = []string{"Severity", "Name", "Method", "API Call", "Created"}
	netRows = append(netRows, kui.netHeaders)
	kui.netTable.Rows = netRows
	kui.netTable.RowSeparator = false
	kui.netTable.Colors.SelectedRowBg = ui.ColorGreen
	kui.netTable.Colors.Text = ui.ColorWhite
	kui.netTable.Colors.SelectedRowBg = ui.ColorBlue
	kui.netTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	kui.netTable.SetRect(1, termHeight/2, termWidth-1, termHeight-1)
	kui.netTable.Title = "K8s API change events"
}

// draw paragraph section
func (kui *KubeKnarkUI) buildParagraph(termWidth, termHeight int) {
	kui.paragraph = widgets.NewParagraph()
	kui.paragraph.Title = "Kube-Knark Tracer"
	kui.paragraph.SetRect(0, 0, termWidth, termHeight)
	kui.paragraph.TextStyle.Fg = ui.ColorWhite
	kui.paragraph.BorderStyle.Fg = ui.ColorCyan
}
