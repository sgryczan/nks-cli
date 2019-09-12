package models

type SolutionTemplateConfig struct {
	Repository     int `json:"repository"`
	RequiredValues map[string]string `json:"required_values"`
	Namespace      string `json:"namespace"`
	Values         string `json:"values_yaml"`
	ChartPath      string `json:"chart_path"`
	Logo           string `json:"logo"`
	ReleaseName    string `json:"release_name"`
	ChartName		string `json:"chart_name"`
	Version			string	`json:"version"`
}

type SolutionTemplateSpec struct {
	Requirements SolutionTemplateSpecRequirements `json:"requirements"`
}

type SolutionTemplateSpecRequirements struct {
	Node 	SolutionTemplateSpecRequirementsNode `json:"node"`
	App     SolutionTemplateSpecRequirementsApp `json:"app"`
}

type SolutionTemplateSpecRequirementsNode struct {
	Count	int	`json:"count"`
	CPU	int	`json:"CPU"`
}
type SolutionTemplateSpecRequirementsApp struct {
	RBAC bool	`json:"rbac"`
	Namespace string `json:"namespace"`
	ValuesEditRequired []string `json:"valuesEditRequired"`
	Storage	bool	`json:"storage"`
	LoadBalancer	bool	`json:"loadbalancer"`
}

type SolutionTemplateDependencies struct {
	Name string `json:"name"`
	Value string `json:"value"`
	Available bool `json:"available"`
	KeysetRequired	bool	`json:"keysetRequired"`
	Tag 	string	`json:"tag"`
	IsPostBuildCompatible bool	`json:"isPostBuildCompatible"`
	IsManagedIndependently bool	`json:"isManagedIndependently"`
	Dependencies	[]string	`json:"dependencies"`
}

type SolutionTemplate struct {
	Name        string         `json:"name"`
	Solution    string         `json:"solution"`
	Installer   string         `json:"installer"`
	Keyset      *int         `json:"keyset"`
	Mode	*string `json:"mode"`
	Tag		string 	`json:"tag"`
	Version     string         `json:"version"`
	Config      SolutionTemplateConfig `json:"config"`
	Spec		SolutionTemplateSpec `json:"spec"`
	Dependencies	SolutionTemplateDependencies `json:"dependencies"`
}