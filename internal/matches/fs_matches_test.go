package matches

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestFSMatches_Full_Match(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := routes.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	routes.BuildMatchMap(fsMap, spec)
	ok := NewFSMatches(fsMap).Match([]string{"chmod", "abc", "/etc/kubernetes/manifests/kube-apiserver.yaml"})
	assert.True(t, ok)
}
func TestFSMatches_Partial_Match(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := routes.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	routes.BuildMatchMap(fsMap, spec)
	ok := NewFSMatches(fsMap).Match([]string{"chmod", "abc", "ddd"})
	assert.False(t, ok)
}

func TestFSMatches_No_Match(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := routes.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	routes.BuildMatchMap(fsMap, spec)
	ok := NewFSMatches(fsMap).Match([]string{"kkk", "abc", "ddd"})
	assert.False(t, ok)
}

func TestFSMatches_EmptyCmd(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := routes.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	routes.BuildMatchMap(fsMap, spec)
	ok := NewFSMatches(fsMap).Match([]string{})
	assert.False(t, ok)
}
