package cmd

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/compiler/mocks"
	mock2 "github.com/chen-keinan/kube-knark/internal/mocks"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestProvideCompiledFiles(t *testing.T) {
	err := utils.CreateKubeKnarkFolders()
	assert.NoError(t, err)
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
	assert.True(t, len(sr[0]) > 0)
	assert.True(t, len(sr[1]) > 0)
}
func TestMatchCmdChan(t *testing.T) {
	kpc := matchCmdChan()
	go func() {
		kpc <- &execevent.KprobeEvent{Pid: uint32(1)}
	}()
	dt := <-kpc
	assert.Equal(t, dt.Pid, uint32(1))
}

func TestMatchNetChan(t *testing.T) {
	kpc := matchNetChan()
	go func() {
		kpc <- &netevent.HTTPNetData{HTTPRequestData: &netevent.HTTPRequestData{Method: common.GET}}
	}()
	dt := <-kpc
	assert.Equal(t, dt.HTTPRequestData.Method, common.GET)
}

func TestNumOfWorkers(t *testing.T) {
	assert.Equal(t, numOfWorkers(), 15)
}

//TestprovideSpecMap provide spec api cache for endpoint validation
func TestProvideSpecMap(t *testing.T) {
	sf := provideSpecFiles()
	sr := provideAPISpecMap(sf)
	fmt.Println(sr)
	assert.True(t, len(sr) > 0)
}

func TestProvideFSSpecMap(t *testing.T) {
	spm := provideFSSpecMap()
	val, ok := spm["chmod"]
	assert.True(t, ok)
	_, ok2 := val.(map[string]interface{})["/etc/kubernetes/manifests/kube-apiserver.yaml"]
	assert.True(t, ok2)
}
