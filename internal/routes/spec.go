package routes

//Spec data model
type Spec struct {
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
	URI         string `mapstructure:"uri" json:"uri"`
	Method      string `mapstructure:"method" json:"uri"`
	Severity    string `mapstructure:"severity" json:"severity"`
}

func (s *Spec) Routes() Routes {
	r := make(Routes, 0)
	for _, c := range s.Categories {
		for _, a := range c.SubCategory.API {
			r = append(r, Route{Method: a.Method, Pattern: a.URI})
		}
	}
	return r
}
