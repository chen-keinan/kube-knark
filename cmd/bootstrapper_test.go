package cmd

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/compiler/mocks"
	mock2 "github.com/chen-keinan/kube-knark/internal/mocks"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestProvideCompiledFiles(t *testing.T) {
	fileData, err := startup.GenerateEbpfFiles()
	assert.NoError(t, err)
	err = startup.SaveFilesIfNotExist(fileData, utils.GetEbpfSourceFolder)
	assert.NoError(t, err)
	ctl := gomock.NewController(t)
	exec := mocks.NewMockExecutor(ctl)
	clang := mock2.NewMockClangExecutor(ctl)
	folder, err := utils.GetEbpfSourceFolder()
	ebpfSourceFolder, err := utils.GetEbpfSourceFolder()
	assert.NoError(t, err)
	ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder()
	assert.NoError(t, err)
	for _, f := range fileData {
		if !strings.HasSuffix(f.Name, ".c") {
			continue
		}
		sourcefilePath := filepath.Join(ebpfSourceFolder, f.Name)
		compiledFilePath := filepath.Join(ebpfCompiledFolder, strings.Replace(f.Name, ".c", ".elf", -1))
		clang.EXPECT().NewExecCommand(fmt.Sprintf(".%s", ebpfSourceFolder), sourcefilePath, compiledFilePath).Return(exec).Times(1)
		clang.EXPECT().CompileSourceToElf(exec).Return(&shell.CommandResult{}, nil).Times(1)
		sf := ProvideCompiledFiles(clang, folder)
		assert.Equal(t, sf[0].Name, "bpf.h")
		assert.Equal(t, sf[1].Name, "bpf_helpers.h")
		assert.Equal(t, sf[2].Name, "kprobe.c")
	}
}

func TestProvideSpecFiles(t *testing.T) {
	sf := ProvideSpecFiles()
	assert.True(t, len(sf[0]) > 0)

}
func TestProvideSpecRoutes(t *testing.T) {
	sf := ProvideSpecFiles()
	sr := ProvideSpecRoutes(sf)
	assert.Equal(t, sr[0][0].Method, common.POST)
	assert.Equal(t, sr[0][0].Pattern, "/api/v1/namespaces/{namespace}/pods")
	assert.Equal(t, sr[0][1].Method, common.PUT)
	assert.Equal(t, sr[0][1].Pattern, "/api/v1/namespaces/{namespace}/pods")
}
func TestMatchCmdChan(t *testing.T) {
	kpc := MatchCmdChan()
	go func() {
		kpc <- &events.KprobeEvent{Pid: uint32(1)}
	}()
	dt := <-kpc
	assert.Equal(t, dt.Pid, uint32(1))
}

func TestMatchNetChan(t *testing.T) {
	kpc := MatchNetChan()
	go func() {
		kpc <- &khttp.HTTPNetData{HttpRequestData: &khttp.HTTPRequestData{Method: common.GET}}
	}()
	dt := <-kpc
	assert.Equal(t, dt.HttpRequestData.Method, common.GET)
}

func TestNumOfWorkers(t *testing.T) {
	assert.Equal(t, NumOfWorkers(), 5)
}
