package ui

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	ui "github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortFsRows(t *testing.T) {
	prevents := []*FilesystemEvt{
		{Spec: &routes.FS{Severity: "MAJOR", SeverityInt: 2}, Msg: &events.KprobeEvent{Args: []string{"a"}}},
		{Spec: &routes.FS{Severity: "CRITICAL", SeverityInt: 1}, Msg: &events.KprobeEvent{Args: []string{"b"}}},
		{Spec: &routes.FS{Severity: "MINOR", SeverityInt: 3}, Msg: &events.KprobeEvent{Args: []string{"c"}}},
	}
	fse := &FilesystemEvt{Spec: &routes.FS{Severity: "INFO", SeverityInt: 4}, Msg: &events.KprobeEvent{Args: []string{"d"}}}
	//prevents = append(prevents,&fse)
	rows := [][]string{{"Severity", "Name", "Command args", "Created"}}
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	sortedRows := nku.sortFSRows(&prevents, fse, rows)
	fmt.Println(len(prevents))
	assert.Equal(t, prevents[0].Spec.SeverityInt, 1)
	assert.Equal(t, prevents[1].Spec.SeverityInt, 2)
	assert.Equal(t, prevents[2].Spec.SeverityInt, 3)
	assert.Equal(t, prevents[3].Spec.SeverityInt, 4)

	assert.Equal(t, sortedRows[1][0], "CRITICAL")
	assert.Equal(t, sortedRows[2][0], "MAJOR")
	assert.Equal(t, sortedRows[3][0], "MINOR")
	assert.Equal(t, sortedRows[4][0], "INFO")
}
func TestSortNetRows(t *testing.T) {
	prevents := []*NetEvt{
		{Spec: &routes.API{Severity: "MAJOR", SeverityInt: 2}, Msg: &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: "GET"}}},
		{Spec: &routes.API{Severity: "CRITICAL", SeverityInt: 1}, Msg: &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: "GET"}}},
		{Spec: &routes.API{Severity: "MINOR", SeverityInt: 3}, Msg: &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: "GET"}}},
	}
	fse := &NetEvt{Spec: &routes.API{Severity: "INFO", SeverityInt: 4}, Msg: &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: "GET"}}}
	//prevents = append(prevents,&fse)
	rows := [][]string{{"Severity", "Name", "Command args", "Created"}}
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	sortedRows := nku.sortNetRows(&prevents, fse, rows)
	fmt.Println(len(prevents))
	assert.Equal(t, prevents[0].Spec.SeverityInt, 1)
	assert.Equal(t, prevents[1].Spec.SeverityInt, 2)
	assert.Equal(t, prevents[2].Spec.SeverityInt, 3)
	assert.Equal(t, prevents[3].Spec.SeverityInt, 4)

	assert.Equal(t, sortedRows[1][0], "CRITICAL")
	assert.Equal(t, sortedRows[2][0], "MAJOR")
	assert.Equal(t, sortedRows[3][0], "MINOR")
	assert.Equal(t, sortedRows[4][0], "INFO")
}

func TestBuildParagraph(t *testing.T) {
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	nku.buildParagraph(100, 200)
	assert.Equal(t, nku.paragraph.Block.Dy(), 200)
	assert.Equal(t, nku.paragraph.Block.Dx(), 100)
	assert.Equal(t, nku.paragraph.TextStyle.Fg, ui.ColorWhite)
	assert.Equal(t, nku.paragraph.BorderStyle.Fg, ui.ColorCyan)
}

func TestNetTable(t *testing.T) {
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	nku.buildNetTable(100, 200)
	assert.Equal(t, nku.netTable.Block.Dy(), 99)
	assert.Equal(t, nku.netTable.Block.Dx(), 98)
	assert.Equal(t, nku.netTable.TextStyle.Fg, ui.ColorWhite)
	assert.Equal(t, nku.netTable.Title, "K8s API change events")
	assert.Equal(t, nku.netHeaders[0], "Severity")
	assert.Equal(t, nku.netHeaders[1], "Name")
	assert.Equal(t, nku.netHeaders[2], "Method")
	assert.Equal(t, nku.netHeaders[3], "API Call")
	assert.Equal(t, nku.netHeaders[4], "Created")
}
func TestFsTable(t *testing.T) {
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	nku.buildFileSystemTable(100, 200)
	assert.Equal(t, nku.fsTable.Block.Dy(), 99)
	assert.Equal(t, nku.fsTable.Block.Dx(), 98)
	assert.Equal(t, nku.fsTable.TextStyle.Fg, ui.ColorWhite)
	assert.Equal(t, nku.fsTable.Title, "K8s configuration file change events")
	assert.Equal(t, nku.fsHeaders[0], "Severity")
	assert.Equal(t, nku.fsHeaders[1], "Name")
	assert.Equal(t, nku.fsHeaders[2], "Command args")
	assert.Equal(t, nku.fsHeaders[3], "Created")
}

func TestWatchEvents(t *testing.T) {
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	uiEvents := make(chan ui.Event)
	err := ui.Init()
	assert.NoError(t, err)
	defer ui.Close()
	// draw external paragraph
	termWidth, termHeight := ui.TerminalDimensions()
	nku.buildParagraph(termWidth, termHeight)
	// init event tables
	nku.buildFileSystemTable(termWidth, termHeight)
	nku.buildNetTable(termWidth, termHeight)
	// render to ui
	ui.Render(nku.paragraph, nku.fsTable, nku.netTable)
	nku.fsTable.Rows = [][]string{
		{"header1", "header2", "header3"},
		{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		{"2016", "10", "11"},
		{"2016", "10", "11"},
		{"2016", "10", "11"}}
	nku.netTable.Rows = [][]string{
		{"header1", "header2", "header3"},
		{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		{"2016", "10", "11"},
		{"2016", "10", "11"},
		{"2016", "10", "11"}}
	go nku.watchEvents(uiEvents)
	uiEvents <- ui.Event{ID: "j"}
	assert.Equal(t, nku.fsTable.curr, nku.fsTable.prev)
	uiEvents <- ui.Event{ID: "k"}
	assert.Equal(t, nku.fsTable.curr, nku.fsTable.prev)
	uiEvents <- ui.Event{ID: "w"}
	assert.Equal(t, nku.netTable.curr, nku.netTable.prev)
	uiEvents <- ui.Event{ID: "s"}
	uiEvents <- ui.Event{ID: "<C-c>"}
}
func TestNewNetEvtChan(t *testing.T) {
	c := NewNetEvtChan()
	go func() {
		c <- NetEvt{Msg: &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: "GET"}}}
	}()
	msg := <-c
	assert.Equal(t, msg.Msg.HTTPRequestData.Method, "GET")
}
func TestNewFilesystemEvtChan(t *testing.T) {
	c := NewFilesystemEvtChan()
	go func() {
		c <- FilesystemEvt{Msg: &events.KprobeEvent{Args: []string{"a"}}}
	}()
	msg := <-c
	assert.Equal(t, msg.Msg.Args[0], "a")
}
