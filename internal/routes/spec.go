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
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URI         string `yaml:"uri"`
	Method      string `yaml:"method"`
	Severity    string `yaml:"severity"`
}

//Routes build routes
func (s *Spec) Routes() Routes {
	r := make(Routes, 0)
	for _, c := range s.Categories {
		for _, a := range c.SubCategory.API {
			r = append(r, Route{Method: a.Method, Pattern: a.URI})
		}
	}
	return r
}
