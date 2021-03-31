package workers

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/kplugin"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
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
	pmd := NewPacketMatchData(rm, pmc, vc, 1, netChan, kplugin.K8sAPICallHook{Plugins: make([]plugin.Symbol, 0)})
	log, err := zap.NewProduction()
	assert.NoError(t, err)
	pmw := NewPacketMatchesWorker(pmd, log)
	pmw.Invoke()
	pmc <- &netevent.HTTPNetData{HTTPRequestData: &netevent.HTTPRequestData{Method: "POST", RequestURI: "/api/v1/namespaces/{namespace}/pods", StartTime: time.Now().String()}}
	msg := <-netChan
	assert.Equal(t, msg.Spec.Severity, "MAJOR")
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
