package routes

import (
	"gopkg.in/yaml.v2"
)

//Routes api routes
type Routes []Route

// A Route defines the parameters for an api endpoint
type Route struct {
	Method  string
	Pattern string
}

//BuildSpecRoutes build api routes from spec file
func BuildSpecRoutes(files []string) ([]Routes, error) {
	routes := make([]Routes, 0)
	for _, f := range files {
		spec := SpecAPI{}
		err := yaml.Unmarshal([]byte(f), &spec)
		if err != nil {
			return nil, err
		}
		routes = append(routes, spec.Routes())
	}
	return routes, nil
}
