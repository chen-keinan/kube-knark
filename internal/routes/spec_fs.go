package routes

import "gopkg.in/yaml.v2"

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
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Commands    []string `yaml:"commands"`
	Severity    string   `yaml:"severity"`
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
	SubFSSpecMap := make(map[string]interface{})
	if len(comm) >= 1 {
		FSSpecMap[comm[0]] = SubFSSpecMap
		buildRecursiveMap(SubFSSpecMap, comm[1:])
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
