package cmd

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/internal/compiler/mocks"
	"github.com/chen-keinan/kube-knark/internal/kplugin"
	mock2 "github.com/chen-keinan/kube-knark/internal/mocks"
	"github.com/chen-keinan/kube-knark/internal/startup"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestProvideCompiledFiles(t *testing.T) {
	fm := utils.NewKFolder()
	err := utils.CreateKubeKnarkFolders(fm)
	assert.NoError(t, err)
	fileData, err := startup.GenerateEbpfFiles()
	assert.NoError(t, err)
	err = startup.SaveFilesIfNotExist(fileData, utils.GetEbpfSourceFolder)
	assert.NoError(t, err)
	ctl := gomock.NewController(t)
	exec := mocks.NewMockExecutor(ctl)
	clang := mock2.NewMockClangExecutor(ctl)
	folder, err := utils.GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	ebpfSourceFolder, err := utils.GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder(fm)
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

func TestProvideFSSpecCache(t *testing.T) {
	spm := provideFSSpecCache()
	val, ok := spm["chmod_*/cni/*_"]
	assert.True(t, ok)
	assert.Equal(t, val.Commands[0], "chmod")
	assert.Equal(t, val.Commands[1], "*/cni/*")
}

func TestLoadAPICallPluginSymbols(t *testing.T) {
	_, err := PluginSetUp("on_k8s_api_call_hook.go")
	assert.NoError(t, err)
	utils.CreateKubeKnarkFolders(utils.NewKFolder())
	sym := LoadAPICallPluginSymbols(NewZapLogger())
	assert.True(t, len(sym.Plugins) > 0)
}

func PluginSetUp(fileName string) (*kplugin.PluginLoader, error) {
	fm := utils.NewKFolder()
	folder, err := utils.GetPluginSourceSubFolder(fm)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(folder)
	if err != nil {
		return nil, err
	}
	cfolder, err := utils.GetCompilePluginSubFolder(fm)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(cfolder)
	if err != nil {
		return nil, err
	}
	err = utils.CreateHomeFolderIfNotExist(fm)
	if err != nil {
		return nil, err
	}
	err = utils.CreatePluginsSourceFolderIfNotExist(fm)
	if err != nil {
		return nil, err
	}
	err = utils.CreatePluginsCompiledFolderIfNotExist(fm)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fmt.Sprintf("./../internal/kplugin/fixtures/%s", fileName))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	nf, err := os.Create(path.Join(folder, "test_plugin.go"))
	if err != nil {
		return nil, err
	}
	_, err = nf.WriteString(string(data))
	if err != nil {
		return nil, err
	}
	pl, err := kplugin.NewPluginLoader()
	if err != nil {
		return nil, err
	}
	return pl, err
}
