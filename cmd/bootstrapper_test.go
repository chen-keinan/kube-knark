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
	assert.NoError(t, err)
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
		sf := provideCompiledFiles(clang, folder)
		assert.Equal(t, sf[0].Name, "bpf.h")
		assert.Equal(t, sf[1].Name, "bpf_helpers.h")
		assert.Equal(t, sf[2].Name, "kprobe.c")
	}
}

func TestProvideSpecFiles(t *testing.T) {
	sf := provideSpecFiles()
	assert.True(t, len(sf[0]) > 0)

}
func TestProvideSpecRoutes(t *testing.T) {
	sf := provideSpecFiles()
	sr := provideSpecRoutes(sf)
	assert.Equal(t, sr[0][0].Method, common.POST)
	assert.Equal(t, sr[0][0].Pattern, "/api/v1/namespaces/{namespace}/pods")
	assert.Equal(t, sr[0][1].Method, common.PUT)
	assert.Equal(t, sr[0][1].Pattern, "/api/v1/namespaces/{namespace}/pods")
}
func TestMatchCmdChan(t *testing.T) {
	kpc := matchCmdChan()
	go func() {
		kpc <- &events.KprobeEvent{Pid: uint32(1)}
	}()
	dt := <-kpc
	assert.Equal(t, dt.Pid, uint32(1))
}

func TestMatchNetChan(t *testing.T) {
	kpc := matchNetChan()
	go func() {
		kpc <- &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: common.GET}}
	}()
	dt := <-kpc
	assert.Equal(t, dt.HTTPRequestData.Method, common.GET)
}

func TestNumOfWorkers(t *testing.T) {
	assert.Equal(t, numOfWorkers(), 5)
}

//TestprovideSpecMap provide spec api cache for endpoint validation
func TestProvideSpecMap(t *testing.T) {
	sf := provideSpecFiles()
	sr := provideAPISpecMap(sf)
	fmt.Println(sr)
	assert.Equal(t, sr[fmt.Sprintf("%s_%s", "POST", "/api/v1/namespaces/{namespace}/pods")].Method, "POST")
	assert.Equal(t, sr[fmt.Sprintf("%s_%s", "POST", "/api/v1/namespaces/{namespace}/pods")].URI, "/api/v1/namespaces/{namespace}/pods")
}

func TestProvideFSSpecMap(t *testing.T) {
	spm := provideFSSpecMap()
	val, ok := spm["chmod"]
	assert.True(t, ok)
	_, ok2 := val.(map[string]interface{})["/etc/kubernetes/manifests/kube-apiserver.yaml"]
	assert.True(t, ok2)
}
