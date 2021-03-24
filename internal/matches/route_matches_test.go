package matches

import (
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteMatches_MatchOK(t *testing.T) {
	a := routes.Routes{
		{
			Method:  common.GET,
			Pattern: "/a/b/{first}/c/{last}",
		},
	}
	router := mux.NewRouter()
	rm := NewRouteMatches([]routes.Routes{a}, router)
	match, err := rm.Match("/a/b/first/c/last", common.GET)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestRouteMatches_MatchBAD(t *testing.T) {
	a := routes.Routes{
		{
			Method:  common.GET,
			Pattern: "/a/b/{first}/c/{last}",
		},
	}
	router := mux.NewRouter()
	rm := NewRouteMatches([]routes.Routes{a}, router)
	match, err := rm.Match("/a/first/c/last", common.GET)
	assert.NoError(t, err)
	assert.True(t, !match)
}
