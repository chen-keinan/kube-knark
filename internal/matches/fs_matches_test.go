package matches

import (
	"bytes"
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
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
	spec := specs.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	specs.BuildMatchMap(fsMap, spec)
	cache, err := specs.CreateFSCacheFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	var sb = new(bytes.Buffer)
	ok := NewFSMatches(fsMap, cache).Match([]string{"chmod", "abc", "/etc/kubernetes/manifests/kube-apiserver.yaml"}, sb)
	assert.Equal(t, sb.String(), "chmod_/etc/kubernetes/manifests/kube-apiserver.yaml_")
	assert.True(t, ok)
	var sb2 = new(bytes.Buffer)
	ok = NewFSMatches(fsMap, cache).Match([]string{"chown", "abc", "/etc/kubernetes/manifests/proxy-apiserver.yaml"}, sb2)
	assert.Equal(t, sb2.String(), "chown_/etc/kubernetes/manifests/proxy-apiserver.yaml_")
	assert.True(t, ok)
}
func TestFSMatches_Partial_Match(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := specs.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	specs.BuildMatchMap(fsMap, spec)
	cache, err := specs.CreateFSCacheFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	var sb = new(bytes.Buffer)
	ok := NewFSMatches(fsMap, cache).Match([]string{"chmod", "abc", "ddd"}, sb)
	assert.Equal(t, sb.String(), "chmod_")
	assert.False(t, ok)
}

func TestFSMatches_No_Match(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := specs.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	specs.BuildMatchMap(fsMap, spec)
	cache, err := specs.CreateFSCacheFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	var sb = new(bytes.Buffer)
	ok := NewFSMatches(fsMap, cache).Match([]string{"kkk", "abc", "ddd"}, sb)
	assert.False(t, ok)
}
func TestFSMatches_Diff_Order(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := specs.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	specs.BuildMatchMap(fsMap, spec)
	cache, err := specs.CreateFSCacheFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	var sb = new(bytes.Buffer)
	ok := NewFSMatches(fsMap, cache).Match([]string{"chmod-rrr", "chmod", "abc", "/etc/kubernetes/manifests/kube-apiserver.yaml"}, sb)
	assert.Equal(t, sb.String(), "chmod_/etc/kubernetes/manifests/kube-apiserver.yaml_")
	assert.True(t, ok)
}

func TestFSMatches_EmptyCmd(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.ConfigFilesPermission))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := specs.SpecFS{}
	yaml.Unmarshal(data, &spec)
	fsMap := make(map[string]interface{})
	specs.BuildMatchMap(fsMap, spec)
	cache, err := specs.CreateFSCacheFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	var sb = new(bytes.Buffer)
	ok := NewFSMatches(fsMap, cache).Match([]string{}, sb)
	assert.False(t, ok)
}
