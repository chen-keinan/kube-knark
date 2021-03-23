package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/internal/tracer/khttp"
)

//PacketMatches instance which match packet data to specific pattern
type PacketMatches struct {
	numOfWorkers int
	pmc          chan *khttp.HTTPNetData
	rm           *matches.RouteMatches
}

//NewMatches return new packet instance
func NewPacketMatches(numOfWorkers int, pmc chan *khttp.HTTPNetData,rm *matches.RouteMatches) *PacketMatches {
	return &PacketMatches{numOfWorkers: numOfWorkers, pmc: pmc,rm:rm}
}

//Invoke invoke packet matches workers
func (pm *PacketMatches) Invoke() {
	for i := 0; i < pm.numOfWorkers; i++ {
		go func() {
			for k := range pm.pmc {
				// display process execution event
				if k.HttpRequestData.RequestURI == "/health" || k.HttpRequestData.RequestURI == "/ready" {
					continue
				}
				kwriter := new(bytes.Buffer)
				err := json.NewEncoder(kwriter).Encode(&k)
				if err != nil {
					continue
				}
				fmt.Println(kwriter.String())
			}
		}()
	}
}
