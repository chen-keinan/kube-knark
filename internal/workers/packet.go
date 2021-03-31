package workers

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
)

//PacketMatchesWorker instance which match packet data to specific pattern
type PacketMatchesWorker struct {
	pmd *PacketMatchData
}

//PacketMatchData encapsulate packet worker properties
type PacketMatchData struct {
	rm              *matches.RouteMatches
	pmc             chan *netevent.HTTPNetData
	validationCache map[string]*specs.API
	uiChan          chan model.NetEvt
	numOfWorkers    int
}

//NewPacketMatchesWorker return new packet instance
func NewPacketMatchesWorker(pmd *PacketMatchData) *PacketMatchesWorker {
	return &PacketMatchesWorker{pmd: pmd}
}

//NewPacketMatchData return new packet data
func NewPacketMatchData(rm *matches.RouteMatches, pmc chan *netevent.HTTPNetData, validationCache map[string]*specs.API, numOfWorkers int, uichan chan model.NetEvt) *PacketMatchData {
	return &PacketMatchData{rm: rm, pmc: pmc, validationCache: validationCache, numOfWorkers: numOfWorkers, uiChan: uichan}
}

//Invoke invoke packet matches workers
func (pm *PacketMatchesWorker) Invoke() {
	for i := 0; i < pm.pmd.numOfWorkers; i++ {
		go func() {
			for k := range pm.pmd.pmc {
				// display process execution event
				if ok, template := pm.pmd.rm.Match(k.HTTPRequestData.RequestURI, k.HTTPRequestData.Method); ok {
					spec := pm.pmd.validationCache[fmt.Sprintf("%s_%s", k.HTTPRequestData.Method, template)]
					pm.pmd.uiChan <- model.NetEvt{Msg: k, Spec: spec}
				}
			}
		}()
	}
}
