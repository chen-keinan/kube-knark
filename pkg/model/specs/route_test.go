package specs

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestBuildSpecRoutes(t *testing.T) {
	f, err := os.Open(fmt.Sprintf("../fixtures/%s", common.Workload))
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	spec := SpecAPI{}
	err = yaml.Unmarshal(data, &spec)
	assert.NoError(t, err)
	routes := spec.Routes()
	assert.Equal(t, routes[0].Method, common.POST)
	assert.Equal(t, routes[1].Method, common.PUT)
	assert.Equal(t, routes[2].Method, common.GET)

}

func TestBuildSpecRoutesError(t *testing.T) {
	data, err := ioutil.ReadAll(strings.NewReader("aaa"))
	assert.NoError(t, err)
	_, err = BuildSpecRoutes([]string{string(data)})
	assert.Error(t, err)
}
