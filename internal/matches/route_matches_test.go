package matches

import (
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

type BC struct {
}

func (c *BC) Routes() Routes {
	return Routes{
		{
			Method:  common.GET,
			Pattern: "/a/b/{first}/c/{last}",
		},
	}
}
func TestRouteMatches_MatchOK(t *testing.T) {
	router := mux.NewRouter()
	rm := NewRouteMatches([]Router{&BC{}}, router)
	match, err := rm.Match("/a/b/chen/c/keinan", common.GET)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestRouteMatches_MatchBAD(t *testing.T) {
	router := mux.NewRouter()
	rm := NewRouteMatches([]Router{&BC{}}, router)
	match, err := rm.Match("/a/chen/c/keinan", common.GET)
	assert.NoError(t, err)
	assert.True(t, !match)
}
