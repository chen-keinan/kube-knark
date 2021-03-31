package matches

import (
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteMatches_MatchOK(t *testing.T) {
	a := specs.Routes{
		{
			Method:  common.GET,
			Pattern: "/a/b/{first}/c/{last}",
		},
	}
	router := mux.NewRouter()
	rm := NewRouteMatches([]specs.Routes{a}, router)
	match, template := rm.Match("/a/b/first/c/last", common.GET)
	assert.Equal(t, template, "/a/b/{first}/c/{last}")
	assert.True(t, match)
}

func TestRouteMatches_MatchBAD(t *testing.T) {
	a := specs.Routes{
		{
			Method:  common.GET,
			Pattern: "/a/b/{first}/c/{last}",
		},
	}
	router := mux.NewRouter()
	rm := NewRouteMatches([]specs.Routes{a}, router)
	match, template := rm.Match("/a/first/c/last", common.GET)
	assert.Equal(t, template, "")
	assert.True(t, !match)
}
