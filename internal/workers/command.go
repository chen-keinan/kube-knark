package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
)

//CommandMatchesWorker instance which match command data to specific pattern
type CommandMatchesWorker struct {
	cmd *CommandMatchData
}

//NewCommandMatchesWorker return new command instance
func NewCommandMatchesWorker(commandMatchData *CommandMatchData) *CommandMatchesWorker {
	return &CommandMatchesWorker{cmd: commandMatchData}
}

//NewCommandMatchesData return new command instance
func NewCommandMatchesData(cmc chan *events.KprobeEvent, NumOfWorkers int, fsMatches *matches.FSMatches) *CommandMatchData {
	return &CommandMatchData{cmc: cmc, numOfWorkers: NumOfWorkers, fsMatches: fsMatches}
}

//CommandMatchData encapsulate command worker properties
type CommandMatchData struct {
	cmc          chan *events.KprobeEvent
	numOfWorkers int
	fsMatches    *matches.FSMatches
}

//Invoke invoke packet matches workers
func (pm *CommandMatchesWorker) Invoke() {
	for i := 0; i < pm.cmd.numOfWorkers; i++ {
		go func() {
			for ke := range pm.cmd.cmc {
				// display process execution event
				if ok := pm.cmd.fsMatches.Match(ke.Args); ok {
					kwriter := new(bytes.Buffer)
					err := json.NewEncoder(kwriter).Encode(&ke)
					if err != nil {
						continue
					}
					fmt.Println(kwriter.String())
				}
			}
		}()
	}
}
