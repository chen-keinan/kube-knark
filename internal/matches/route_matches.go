package matches

import (
	"github.com/gorilla/mux"
	"net/http"
)

//Routes api routes
type Routes []Route

// A Route defines the parameters for an api endpoint
type Route struct {
	Method  string
	Pattern string
}

//Router interface include all routes pattern and method
type Router interface {
	Routes() Routes
}

//RouteMatches Object
type RouteMatches struct {
	router *mux.Router
}

//NewRouteMatches create new routes matches instance
func NewRouteMatches(ValidationRoutes []Router, router *mux.Router) *RouteMatches {
	for _, mo := range ValidationRoutes {
		for _, rt := range mo.Routes() {
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
