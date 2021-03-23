package matches

import (
	"github.com/chen-keinan/kube-knark/internal/routes"
	"github.com/gorilla/mux"
	"net/http"
)

//RouteMatches Object
type RouteMatches struct {
	router *mux.Router
}

//NewRouteMatches create new routes matches instance
func NewRouteMatches(ValidationRoutes []routes.Routes, router *mux.Router) *RouteMatches {
	for _, mo := range ValidationRoutes {
		for _, rt := range mo {
			router.Methods(rt.Method).Path(rt.Pattern)
		}
	}
	return &RouteMatches{router: router}
}

//Match match route to specified api pattern
func (mr RouteMatches) Match(url, method string) (bool, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, err
	}
	return mr.router.Match(request, &mux.RouteMatch{}), nil
}
