package specs

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
	sMap := make(map[string]interface{})
	BuildMatchMap(sMap, spec)
	val, ok := sMap["chmod"]
	assert.True(t, ok)
	_, ok2 := val.(map[string]interface{})["/etc/kubernetes/manifests/kube-apiserver.yaml"]
	assert.True(t, ok2)
}

func TestCreateFSMapFromSpecFiles(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	fsm, err := CreateFSMapFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	val, ok := fsm["chmod"]
	assert.True(t, ok)
	_, ok2 := val.(map[string]interface{})["/etc/kubernetes/manifests/kube-apiserver.yaml"]
	_, ok3 := val.(map[string]interface{})["/etc/kubernetes/manifests/proxy-apiserver.yaml"]
	assert.True(t, ok2)
	assert.True(t, ok3)
}

func TestCreateFSCacheFromSpecFiles(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	cache, err := CreateFSCacheFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	assert.True(t, cache["chmod_/etc/kubernetes/manifests/kube-apiserver.yaml_"].Name == "change kube-apiserver.yaml file permission")
	assert.True(t, cache["chmod_/etc/kubernetes/manifests/proxy-apiserver.yaml_"].Name == "change proxy-apiserver.yaml file permission")
}
func TestUnmarshallFS(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := SpecFS{}
	err = yaml.Unmarshal(data, &spec)
	assert.NoError(t, err)
	assert.Equal(t, spec.Categories[0].SubCategory.FS[0].SeverityInt, 1)
}
