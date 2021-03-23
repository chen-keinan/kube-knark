package routes

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/stretchr/testify/assert"
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
	routes, err := BuildSpecRoutes([]string{string(data)})
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, routes[0][0].Method, common.POST)
	assert.Equal(t, routes[0][1].Method, common.PUT)

}

func TestBuildSpecRoutesError(t *testing.T) {
	data, err := ioutil.ReadAll(strings.NewReader("aaa"))
	assert.NoError(t, err)
	_, err = BuildSpecRoutes([]string{string(data)})
	assert.Error(t, err)
}
