package workers

import (
	"bytes"
	"fmt"
	"github.com/chen-keinan/go-user-plugins/uplugin"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"go.uber.org/zap"
	"plugin"
)

//K8sFileConfigChangeHook hold the plugin symbol for K8s File Config Change Hook
type K8sFileConfigChangeHook struct {
	Plugins []plugin.Symbol
	Plug    *uplugin.PluginLoader
}

//CommandMatchesWorker instance which match command data to specific pattern
type CommandMatchesWorker struct {
	cmd *CommandMatchData
	log *zap.Logger
}

//NewCommandMatchesWorker return new command instance
func NewCommandMatchesWorker(commandMatchData *CommandMatchData, log *zap.Logger) *CommandMatchesWorker {
	return &CommandMatchesWorker{cmd: commandMatchData, log: log}
}

//NewCommandMatchesData return new command instance
func NewCommandMatchesData(cmc chan *execevent.KprobeEvent, NumOfWorkers int, fsMatches *matches.FSMatches, uiChan chan model.K8sConfigFileChangeEvent, hook K8sFileConfigChangeHook) *CommandMatchData {
	return &CommandMatchData{cmc: cmc, numOfWorkers: NumOfWorkers, fsMatches: fsMatches, uiChan: uiChan, plugins: hook}
}

//CommandMatchData encapsulate command worker properties
type CommandMatchData struct {
	cmc          chan *execevent.KprobeEvent
	numOfWorkers int
	fsMatches    *matches.FSMatches
	uiChan       chan model.K8sConfigFileChangeEvent
	plugins      K8sFileConfigChangeHook
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
					evt := model.K8sConfigFileChangeEvent{Msg: ke, Spec: fSpec}
					pm.cmd.uiChan <- evt
					if len(pm.cmd.plugins.Plugins) > 0 {
						for _, pl := range pm.cmd.plugins.Plugins {
							_, err := pm.cmd.plugins.Plug.Invoke(pl, evt)
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
