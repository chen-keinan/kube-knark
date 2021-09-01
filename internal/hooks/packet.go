package hooks

import (
	"fmt"
	"github.com/chen-keinan/go-user-plugins/uplugin"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
	"go.uber.org/zap"
	"plugin"
)

//K8sAPICallHook hold the plugin symbol for K8s API Call Hook
type K8sAPICallHook struct {
	Plugins []plugin.Symbol
	Plug    *uplugin.PluginLoader
}

//PacketMatchesWorker instance which match packet data to specific pattern
type PacketMatchesWorker struct {
	pmd *PacketMatchData
	log *zap.Logger
}

//PacketMatchData encapsulate packet worker properties
type PacketMatchData struct {
	rm              *matches.RouteMatches
	pmc             chan *netevent.HTTPNetData
	validationCache map[string]*specs.API
	uiChan          chan model.K8sAPICallEvent
	numOfWorkers    int
	plugins         K8sAPICallHook
}

//NewPacketMatchesWorker return new packet instance
func NewPacketMatchesWorker(pmd *PacketMatchData, log *zap.Logger) *PacketMatchesWorker {
	return &PacketMatchesWorker{pmd: pmd, log: log}
}

//NewPacketMatchData return new packet data
func NewPacketMatchData(rm *matches.RouteMatches, pmc chan *netevent.HTTPNetData, validationCache map[string]*specs.API, numOfWorkers int, uichan chan model.K8sAPICallEvent, plugin K8sAPICallHook) *PacketMatchData {
	return &PacketMatchData{rm: rm, pmc: pmc, validationCache: validationCache, numOfWorkers: numOfWorkers, uiChan: uichan, plugins: plugin}
}

//Invoke invoke packet matches workers
func (pm *PacketMatchesWorker) Invoke() {
	for i := 0; i < pm.pmd.numOfWorkers; i++ {
		go func() {
			for k := range pm.pmd.pmc {
				// display process execution event
				if ok, template := pm.pmd.rm.Match(k.HTTPRequestData.RequestURI, k.HTTPRequestData.Method); ok {
					spec := pm.pmd.validationCache[fmt.Sprintf("%s_%s", k.HTTPRequestData.Method, template)]
					evt := model.K8sAPICallEvent{Msg: k, Spec: spec}
					pm.pmd.uiChan <- evt
					if len(pm.pmd.plugins.Plugins) > 0 {
						for _, pl := range pm.pmd.plugins.Plugins {
							_, err := pm.pmd.plugins.Plug.Invoke(pl, evt)
							if err != nil {
								pm.log.Error(fmt.Sprintf("failed to execute plugins %s", err.Error()))
							}
						}
					}
				}
			}
		}()
	}
}
