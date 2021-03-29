package routes

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"strings"
)

//SpecFS data model
type SpecFS struct {
	SpecFile   string       `yaml:"spec"`
	Categories []CategoryFS `yaml:"categories"`
}

//CategoryFS data model
type CategoryFS struct {
	Name        string         `yaml:"name"`
	SubCategory *SubCategoryFS `yaml:"sub_category"`
}

//SubCategoryFS data model
type SubCategoryFS struct {
	Name string `yaml:"name"`
	FS   []*FS  `yaml:"fs"`
}

//FS data model
type FS struct {
	Name        string   `mapstructure:"name" yaml:"name"`
	Description string   `mapstructure:"description" yaml:"description"`
	Commands    []string `mapstructure:"commands" yaml:"commands"`
	Severity    string   `mapstructure:"severity" yaml:"severity"`
	SeverityInt int      `mapstructure:"severity_int" yaml:"severity_int"`
}

//UnmarshalYAML over unmarshall
func (at *FS) UnmarshalYAML(unmarshal func(interface{}) error) error {
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

//BuildMatchMap build fs match map
func BuildMatchMap(FSSpecMap map[string]interface{}, s SpecFS) {
	for _, c := range s.Categories {
		for _, a := range c.SubCategory.FS {
			buildRecursiveMap(FSSpecMap, a.Commands)
		}
	}
}

//buildRecursiveMap build fs map helper
func buildRecursiveMap(FSSpecMap map[string]interface{}, comm []string) {
	if len(comm) >= 1 {
		subFSSpecMap, ok := FSSpecMap[comm[0]]
		if !ok {
			subFSSpecMap = make(map[string]interface{})
			FSSpecMap[comm[0]] = subFSSpecMap
		}
		buildRecursiveMap(subFSSpecMap.(map[string]interface{}), comm[1:])
	}
}

//CreateFSMapFromSpecFiles build spec api cache for presenting trace data
func CreateFSMapFromSpecFiles(specFiles []string) (map[string]interface{}, error) {
	specMap := make(map[string]interface{})
	for _, f := range specFiles {
		spec := SpecFS{}
		err := yaml.Unmarshal([]byte(f), &spec)
		if err != nil {
			return nil, err
		}
		BuildMatchMap(specMap, spec)
	}
	return specMap, nil
}

//CreateFSCacheFromSpecFiles build spec fs cache for data processing
func CreateFSCacheFromSpecFiles(specFiles []string) (map[string]*FS, error) {
	specMap := make(map[string]*FS)
	for _, f := range specFiles {
		spec := SpecFS{}
		err := yaml.Unmarshal([]byte(f), &spec)
		if err != nil {
			return nil, err
		}
		for _, categories := range spec.Categories {
			for _, fs := range categories.SubCategory.FS {
				var sb strings.Builder
				for _, com := range fs.Commands {
					sb.WriteString(fmt.Sprintf("%s_", com))
				}
				specMap[sb.String()] = fs
			}
		}
	}
	return specMap, nil
}
