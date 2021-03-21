package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
)

//PacketMatches instance which match command data to specific pattern
type CommandMatches struct {
	numOfWorkers int
	cmc          chan *events.KprobeEvent
}

//NewCommandMatches return new command instance
func NewCommandMatches(NumOfWorkers int, cmc chan *events.KprobeEvent) *CommandMatches {
	return &CommandMatches{numOfWorkers: NumOfWorkers, cmc: cmc}
}

//Invoke invoke packet matches workers
func (pm *CommandMatches) Invoke() {
	for i := 0; i < pm.numOfWorkers; i++ {
		go func() {
			for ke := range pm.cmc {
				// display process execution event
				kwriter := new(bytes.Buffer)
				err := json.NewEncoder(kwriter).Encode(&ke)
				if err != nil {
					continue
				}
				fmt.Println(kwriter.String())
			}
		}()
	}
}
