package kplugin

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestPluginLoader_Plugins(t *testing.T) {
	pl, err := pluginSetUp("on_k8s_api_call_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	assert.Equal(t, plFiles[0], "test_plugin.go")
}
func TestExecuteK8sConfigChange(t *testing.T) {
	pl, err := pluginSetUp("on_k8s_file_config_change_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	sym, err := pl.Compile(plFiles[0], common.OnK8sFileConfigChangeHook)
	assert.NoError(t, err)
	err = ExecuteK8sConfigChange(sym, model.K8sConfigFileChangeEvent{})
	assert.NoError(t, err)
}

func TestExecuteNetEvt(t *testing.T) {
	pl, err := pluginSetUp("on_k8s_api_call_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	sym, err := pl.Compile(plFiles[0], common.OnK8sAPICallHook)
	assert.NoError(t, err)
	err = ExecuteNetEvt(sym, model.K8sAPICallEvent{})
	assert.NoError(t, err)
}

func TestPluginLoader_CompileBad(t *testing.T) {
	pl, err := pluginSetUp("empty.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	_, err = pl.Compile(plFiles[0], common.OnK8sAPICallHook)
	assert.Error(t, err)
	_, err = pl.Compile("a/b/c", common.OnK8sAPICallHook)
	assert.Error(t, err)
}
func TestPluginLoader_CompileWrongHook(t *testing.T) {
	pl, err := pluginSetUp("on_k8s_api_call_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	_, err = pl.Compile(plFiles[0], "NoHook")
	assert.Error(t, err)
}

func pluginSetUp(fileName string) (*PluginLoader, error) {
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
	f, err := os.Open(fmt.Sprintf("./fixtures/%s", fileName))
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
	pl, err := NewPluginLoader()
	if err != nil {
		return nil, err
	}
	return pl, err
}
