package versioning

type Module struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Provider  string `json:"provider"`
	Source    string `json:"source"`
	Tag       string `json:"tag"`
	Root	  Root	 `json:"root"`
	Submodules []Submodule`json:"submodules"`
	Providers []string `json:"providers"`
	Versions  []string `json:"versions"`
}

type Root struct{
	Inputs []Input			`json:"inputs"`
	Outputs []Output		`json:"outputs"`
	Resources []Resource	`json:"resources"`
}

type Input struct{
	Name string		`json:"name"`
	Type     string `json:"type"`
	Required bool	`json:"required"`
}

type Output struct{
	Name string `json:"name"`
}

type Resource struct{
	Name string	`json:"name"`
	Type string	`json:"type"`
}

type Submodule struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Inputs []Input `json:"inputs"`
	Outputs []Output `json:"outputs"`
	Resources []Resource `json:"resources"`
}

