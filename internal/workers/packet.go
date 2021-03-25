package workers

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
	"github.com/chen-keinan/kube-knark/pkg/ui"
)

//PacketMatchesWorker instance which match packet data to specific pattern
type PacketMatchesWorker struct {
	pmd *PacketMatchData
}

//PacketMatchData encapsulate packet worker properties
type PacketMatchData struct {
	rm              *matches.RouteMatches
	pmc             chan *khttp.HTTPNetData
	validationCache map[string]*routes.API
	uiChan          chan ui.NetEvt
	numOfWorkers    int
}

//NewPacketMatchesWorker return new packet instance
func NewPacketMatchesWorker(pmd *PacketMatchData) *PacketMatchesWorker {
	return &PacketMatchesWorker{pmd: pmd}
}

//NewPacketMatchData return new packet data
func NewPacketMatchData(rm *matches.RouteMatches, pmc chan *khttp.HTTPNetData, validationCache map[string]*routes.API, numOfWorkers int, uichan chan ui.NetEvt) *PacketMatchData {
	return &PacketMatchData{rm: rm, pmc: pmc, validationCache: validationCache, numOfWorkers: numOfWorkers, uiChan: uichan}
}

//Invoke invoke packet matches workers
func (pm *PacketMatchesWorker) Invoke() {
	for i := 0; i < pm.pmd.numOfWorkers; i++ {
		go func() {
			for k := range pm.pmd.pmc {
				// display process execution event
				if ok, _ := pm.pmd.rm.Match(k.HTTPRequestData.RequestURI, k.HTTPRequestData.Method); ok {
					spec := pm.pmd.validationCache[fmt.Sprintf("%s_%s", k.HTTPRequestData.Method, k.HTTPRequestData.RequestURI)]
					pm.pmd.uiChan <- ui.NetEvt{Msg: k, Spec: spec}
				}
			}
		}()
	}
}
