package routes

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

//TestBuildMatchMap
func TestBuildMatchMap(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := BuildMatchMap(spec)
	val, ok := fsMap["chmod"]
	assert.True(t, ok)
	_, ok2 := val.(map[string]interface{})["/etc/kubernetes/manifests/kube-apiserver.yaml"]
	assert.True(t, ok2)
}
