package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
)

//PacketMatches instance which match packet data to specific pattern
type PacketMatches struct {
	pmd *PacketMatchData
}

//PacketMatchData encapsulate packet worker properties
type PacketMatchData struct {
	rm           *matches.RouteMatches
	pmc          chan *khttp.HTTPNetData
	cache        map[string]*routes.API
	numOfWorkers int
}

//NewPacketMatches return new packet instance
func NewPacketMatches(pmd *PacketMatchData) *PacketMatches {
	return &PacketMatches{pmd: pmd}
}

//NewPacketMatchData return new packet data
func NewPacketMatchData(rm *matches.RouteMatches, pmc chan *khttp.HTTPNetData, cache map[string]*routes.API, numOfWorkers int) *PacketMatchData {
	return &PacketMatchData{rm: rm, pmc: pmc, cache: cache, numOfWorkers: numOfWorkers}
}

//Invoke invoke packet matches workers
func (pm *PacketMatches) Invoke() {
	for i := 0; i < pm.pmd.numOfWorkers; i++ {
		go func() {
			for k := range pm.pmd.pmc {
				// display process execution event
				kwriter := new(bytes.Buffer)
				err := json.NewEncoder(kwriter).Encode(&k)
				if err != nil {
					continue
				}
				if ok, _ := pm.pmd.rm.Match(k.HTTPRequestData.RequestURI, k.HTTPRequestData.Method); ok {
					fmt.Println(kwriter.String())
				}
			}
		}()
	}
}
