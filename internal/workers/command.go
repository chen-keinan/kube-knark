package workers

import (
	"bytes"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/ui"
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
func NewCommandMatchesData(cmc chan *execevent.KprobeEvent, NumOfWorkers int, fsMatches *matches.FSMatches, uiChan chan ui.FilesystemEvt) *CommandMatchData {
	return &CommandMatchData{cmc: cmc, numOfWorkers: NumOfWorkers, fsMatches: fsMatches, uiChan: uiChan}
}

//CommandMatchData encapsulate command worker properties
type CommandMatchData struct {
	cmc          chan *execevent.KprobeEvent
	numOfWorkers int
	fsMatches    *matches.FSMatches
	uiChan       chan ui.FilesystemEvt
}

//Invoke invoke packet matches workers
func (pm *CommandMatchesWorker) Invoke() {
	for i := 0; i < pm.cmd.numOfWorkers; i++ {
		go func() {
			for ke := range pm.cmd.cmc {
				// display process execution event
				var sb = new(bytes.Buffer)
				if ok := pm.cmd.fsMatches.Match(ke.Args, sb); ok {
					fSpec := pm.cmd.fsMatches.Cache[sb.String()]
					pm.cmd.uiChan <- ui.FilesystemEvt{Msg: ke, Spec: fSpec}
				}
			}
		}()
	}
}
