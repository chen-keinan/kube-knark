package workers

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/kplugin"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"testing"
	"time"
)

func TestPacketMatchesWorker_Invoke(t *testing.T) {
	pmc := make(chan *netevent.HTTPNetData)
	sr, err := buildSpecRoutes()
	assert.NoError(t, err)
	rm := matches.NewRouteMatches(sr, mux.NewRouter())
	vc, err := buildValidationCache()
	assert.NoError(t, err)
	netChan := make(chan model.K8sAPICallEvent)
	symbols, err := getPluginSymbols("on_k8s_api_call_hook.go", common.OnK8sAPICallHook)
	assert.NoError(t, err)
	pmd := NewPacketMatchData(rm, pmc, vc, 1, netChan, kplugin.K8sAPICallHook{Plugins: symbols})
	log, err := zap.NewProduction()
	assert.NoError(t, err)
	pmw := NewPacketMatchesWorker(pmd, log)
	assert.NoError(t, err)
	pmw.Invoke()
	pmc <- &netevent.HTTPNetData{HTTPRequestData: &netevent.HTTPRequestData{Method: "POST", RequestURI: "/api/v1/namespaces/{namespace}/pods", StartTime: time.Now().String()}}
	msg := <-netChan
	assert.Equal(t, msg.Spec.Severity, "MAJOR")
}

func getPluginSymbols(name, method string) ([]plugin.Symbol, error) {
	pl, err := pluginSetUp(name)
	if err != nil {
		return nil, err
	}
	plugins, err := pl.Plugins()
	if err != nil {
		return nil, err
	}
	symbols := make([]plugin.Symbol, 0)
	for _, name := range plugins {
		sym, err := pl.Compile(name, method)
		if err != nil {
			return nil, err
		}
		symbols = append(symbols, sym)
	}
	return symbols, nil
}

func buildSpecRoutes() ([]specs.Routes, error) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return specs.BuildSpecRoutes([]string{string(data)})
}

func buildValidationCache() (map[string]*specs.API, error) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return specs.CreateMapFromSpecFiles([]string{string(data)})
}

func pluginSetUp(fileName string) (*kplugin.PluginLoader, error) {
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
	f, err := os.Open(fmt.Sprintf("./../kplugin/fixtures/%s", fileName))
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
