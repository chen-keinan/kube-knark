package specs

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

//SpecAPI data model
type SpecAPI struct {
	SpecFile   string     `yaml:"spec"`
	Categories []Category `yaml:"categories"`
}

//Category data model
type Category struct {
	Name        string       `yaml:"name"`
	SubCategory *SubCategory `yaml:"sub_category"`
}

//SubCategory data model
type SubCategory struct {
	Name string `yaml:"name"`
	API  []*API `yaml:"api"`
}

//API data model
type API struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Description string `mapstructure:"description" yaml:"description"`
	URI         string `mapstructure:"uri" yaml:"uri"`
	Method      string `mapstructure:"method" yaml:"method"`
	Severity    string `mapstructure:"severity" yaml:"severity"`
	SeverityInt int    `mapstructure:"severity_int" yaml:"severity_int"`
}

//UnmarshalYAML over unmarshall
func (at *API) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var res map[string]interface{}
	if err := unmarshal(&res); err != nil {
		return err
	}
	err := mapstructure.Decode(res, &at)
	if err != nil {
		return err
	}
	at.SeverityInt = utils.FindSeverityInt(at.Severity)
	return nil
}

//Routes build routes
func (s *SpecAPI) Routes() Routes {
	r := make(Routes, 0)
	for _, c := range s.Categories {
		for _, a := range c.SubCategory.API {
			r = append(r, Route{Method: a.Method, Pattern: a.URI})
		}
	}
	return r
}

func buildSpecMap(specMap map[string]*API, spec *SpecAPI) {
	for _, s := range spec.Categories {
		for _, a := range s.SubCategory.API {
			specMap[fmt.Sprintf("%s_%s", a.Method, a.URI)] = a
		}
	}
}

//CreateMapFromSpecFiles build spec api cache for presenting trace data
func CreateMapFromSpecFiles(specFiles []string) (map[string]*API, error) {
	specMap := make(map[string]*API)
	for _, f := range specFiles {
		spec := SpecAPI{}
		err := yaml.Unmarshal([]byte(f), &spec)
		if err != nil {
			return nil, err
		}
		buildSpecMap(specMap, &spec)
	}
	return specMap, nil
}
