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

func TestRoutes(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	assert.NoError(t, err)
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	routes, err := BuildSpecRoutes([]string{string(data)})
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, routes[0][0].Method, common.POST)
	assert.Equal(t, routes[0][1].Method, common.PUT)

}
func TestCreateMapFromSpecFiles(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := Spec{}
	yaml.Unmarshal(data, &spec)
	mapAPI, err := CreateMapFromSpecFiles([]string{string(data)})
	assert.NoError(t, err)
	api := spec.Categories[0].SubCategory.API[0]
	assert.Equal(t, mapAPI[fmt.Sprintf("%s_%s", api.Method, api.URI)].Method, api.Method)
	assert.Equal(t, mapAPI[fmt.Sprintf("%s_%s", api.Method, api.URI)].URI, api.URI)
}
