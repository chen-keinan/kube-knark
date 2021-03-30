package workers

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/ui"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestPacketMatchesWorker_Invoke(t *testing.T) {
	pmc := make(chan *khttp.HTTPNetData)
	sr, err := buildSpecRoutes()
	assert.NoError(t, err)
	rm := matches.NewRouteMatches(sr, mux.NewRouter())
	vc, err := buildValidationCache()
	assert.NoError(t, err)
	netChan := make(chan ui.NetEvt)
	pmd := NewPacketMatchData(rm, pmc, vc, 1, netChan)
	pmw := NewPacketMatchesWorker(pmd)
	pmw.Invoke()
	pmc <- &khttp.HTTPNetData{HTTPRequestData: &khttp.HTTPRequestData{Method: "POST", RequestURI: "/api/v1/namespaces/{namespace}/pods", StartTime: time.Now().String()}}
	msg := <-netChan
	assert.Equal(t, msg.Spec.Severity, "MAJOR")
}

func buildSpecRoutes() ([]routes.Routes, error) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return routes.BuildSpecRoutes([]string{string(data)})
}

func buildValidationCache() (map[string]*routes.API, error) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return routes.CreateMapFromSpecFiles([]string{string(data)})
}
