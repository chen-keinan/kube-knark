package workers

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/matches"
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
	"github.com/chen-keinan/kube-knark/pkg/ui"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCommandMatchesWorker_Invoke(t *testing.T) {
	cmc := make(chan *execevent.KprobeEvent)
	mmap, err := buildMatchMap()
	assert.NoError(t, err)
	smap, err := getSpecMap()
	assert.NoError(t, err)
	fsMatches := matches.NewFSMatches(mmap, smap)
	uichan := make(chan ui.FilesystemEvt)
	cmd := NewCommandMatchesData(cmc, 1, fsMatches, uichan)
	cmw := NewCommandMatchesWorker(cmd)
	cmw.Invoke()
	cmc <- &execevent.KprobeEvent{StartTime: time.Now().String(), UID: uint32(1), Pid: uint32(1), Gid: uint32(1), Comm: "cmd", Args: []string{"chmod", "/etc/kubernetes/manifests/kube-apiserver.yaml"}}
	res := <-uichan
	assert.Equal(t, res.Spec.Severity, "CRITICAL")
}

func buildMatchMap() (map[string]interface{}, error) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	spec := specs.SpecFS{}
	yaml.Unmarshal(data, &spec)
	sMap := make(map[string]interface{})
	specs.BuildMatchMap(sMap, spec)
	return sMap, nil
}

func getSpecMap() (map[string]*specs.FS, error) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	fsm, err := specs.CreateFSCacheFromSpecFiles([]string{string(data)})
	if err != nil {
		return nil, err
	}
	return fsm, nil
}
