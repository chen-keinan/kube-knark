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
	p := nku.buildParagraph(100, 200)
	assert.Equal(t, p.Block.Dy(), 200)
	assert.Equal(t, p.Block.Dx(), 100)
	assert.Equal(t, p.TextStyle.Fg, ui.ColorWhite)
	assert.Equal(t, p.BorderStyle.Fg, ui.ColorCyan)
}

func TestNetTable(t *testing.T) {
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	p, headers := nku.buildNetTable(100, 200)
	assert.Equal(t, p.Block.Dy(), 99)
	assert.Equal(t, p.Block.Dx(), 98)
	assert.Equal(t, p.TextStyle.Fg, ui.ColorWhite)
	assert.Equal(t, p.Title, "K8s API change events")
	assert.Equal(t, headers[0], "Severity")
	assert.Equal(t, headers[1], "Name")
	assert.Equal(t, headers[2], "Method")
	assert.Equal(t, headers[3], "API Call")
	assert.Equal(t, headers[4], "Created")
}
func TestFsTable(t *testing.T) {
	nku := NewKubeKnarkUI(make(chan NetEvt), make(chan FilesystemEvt))
	p, headers := nku.buildFileSystemTable(100, 200)
	assert.Equal(t, p.Block.Dy(), 99)
	assert.Equal(t, p.Block.Dx(), 98)
	assert.Equal(t, p.TextStyle.Fg, ui.ColorWhite)
	assert.Equal(t, p.Title, "K8s configuration file change events")
	assert.Equal(t, headers[0], "Severity")
	assert.Equal(t, headers[1], "Name")
	assert.Equal(t, headers[2], "Command args")
	assert.Equal(t, headers[3], "Created")
}
