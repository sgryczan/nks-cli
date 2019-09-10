package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SolutionConfig struct {
	Repository     string `json:"repository"`
	RequiredValues map[string]string
	Namespace      string `json:"namespace"`
	Values         string `json:"values_yaml"`
	ChartPath      string `json:"chart_path"`
	Logo           string `json:"logo"`
	ReleaseName    string `json:"release_name"`
}

type Solution struct {
	ID          int            `json:"pk"`
	Name        string         `json:"name"`
	InstanceID  string         `json:"instance_id"`
	Cluster     int            `json:"cluster"`
	Solution    string         `json:"solution"`
	Installer   string         `json:"installer"`
	Keyset      string         `json:"keyset"`
	KeysetName  string         `json:"keyset_name"`
	Version     string         `json:"version"`
	State       string         `json:"state"`
	URL         string         `json:"url"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	MaxNodes    string         `json:"max_nodes"`
	GitRepo     string         `json:"git_repo"`
	GitPath     string         `json:"git_path"`
	Initial     bool           `json:"initial"`
	Config      SolutionConfig `json:"config"`
	ExtraData   map[string]string
	CreatedTime string `json:"created"`
	UpdatedTime string `json:"updated"`
	IsDeletable bool   `json:"is_deleteable"`
}

// solutionCmd represents the solution command
var getSolutionsCmd = &cobra.Command{
	Use:   "solutions",
	Short: "list solutions",
	Run: func(cmd *cobra.Command, args []string) {
		ss, err := getSolutions()
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			ss = &[]Solution{}
		}
		printSolutions(*ss)
	},
}

func getSolutions() (*[]Solution, error) {
	orgID := viper.GetString("org_id")
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%s/solutions", orgID)
	res, err := httpRequest("GET", url, "")

	data := []Solution{}

	_ = json.Unmarshal(res, &data)
	//check(err)

	return &data, err
}

func printSolutions(s []Solution) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tSOLUTION\tSTATE\t\n")
	for _, s := range s {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", s.Name, s.ID, s.Solution, s.State)
	}
	w.Flush()
}

func init() {
	getCmd.AddCommand(getSolutionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// solutionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// solutionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
