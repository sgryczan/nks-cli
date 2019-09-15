package models

import (
	"encoding/json"
	"text/tabwriter"
	"os"
	"fmt"
)

var solutionTemplates = map[string]SolutionTemplate{
	"jenkins":  jenkins,
}

func GetTemplateAsJson(s string) (string, error) {
	t := solutionTemplates[s]
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("error while attempting to get solution template %s:\n\t:%v", s, err)
	}
	return string(b), err
}

func RepositoryToTemplate(r Repository, releaseName string) *SolutionTemplate {
	template := SolutionTemplate{}

	chartName := r.ChartIndex[0].Name
	chartVersion := r.ChartIndex[0].Chart["version"]
	chartLogo := "/images/k8s-logo-border.ae630e27.png"
	
	template.Name = chartName
	template.Solution = fmt.Sprintf("%s-repo-%d", chartName, r.ID)
	template.Installer = "helm"
	template.Keyset = r.KeysetId
	template.Mode = nil
	template.Tag = chartVersion
	template.Config = SolutionTemplateConfig{
		Namespace: chartName,
		ChartName: chartName,
		Version: chartVersion,
		ChartPath: r.Path,
		ReleaseName: releaseName,
		Logo: chartLogo,
		Repository: r.ID,
		Values: r.ChartIndex[0].Values,
		RequiredValues: map[string]string{},
	}
	template.Spec = SolutionTemplateSpec{}
	template.Dependencies = SolutionTemplateDependencies{
		Name: "Helm Tiller",
		Value: "helm_tiller",
		Available: false,
		KeysetRequired: false,
		Tag: "latest",
		IsPostBuildCompatible: true,
		IsManagedIndependently: false,
		Dependencies: []string{},
	}
	template.Version = chartVersion

	return &template
}

func ListSolutionTemplates() []SolutionTemplate {

	solutions := []SolutionTemplate{}

	for _, val := range solutionTemplates {
		solutions = append(solutions, val)
	}
	return solutions
}


func PrintSolutionTemplates(t *[]SolutionTemplate) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tVERSION\tCHART PATH\tINSTALLER\t\n")
	for _, v := range *t {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.Name, (v.Config).Version, (v.Config).ChartPath, v.Installer)
	}
	w.Flush()
}