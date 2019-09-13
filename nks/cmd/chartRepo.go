package cmd

import (
	"fmt"
	"text/tabwriter"
	"os"
	"strings"
	"encoding/json"

	"github.com/spf13/cobra"
)

var flagCreateRepositorySourceType string

type createRepoInput []checkRepositoryInput

type createRepoResponseS []createRepoResponse

type createRepoResponse struct {
	ID          int         `json:"pk"`
	Name        string    	`json:"name"`
	Source		string		`json:"source"`
	SourceDisplay	string 	`json:"Github"`
	OrganizationId		int 	`json:"org"`
	Path	string			`json:"path"`
	URL		string			`json:"url"`
	IsSystem	bool		`json:"is_system"`
	IsPrivate bool			`json:"is_private"`
	KeysetId	*int		`json:"keyset"`
	ChartIndex	[]chartIndex	`json:"chart_index"`
	State      string		`json:"state"`
	Owner		int		`json:"owner"`
	IsAccessible	bool	`json:"is_accessible"`
	Synced		*string		`json:"synced"`
	Created 	string		`json:"created"`
	Updated 	string		`json:"updated"`
}

type chartIndex struct {
	Name        string    	`json:"name"`
	Sha			  string		`json:"sha"`
	Chart		map[string]string `json:"chart"`
	Values		string			`json:"values"`
	Path		string		`json:"path"`
	Spec		map[string]string `json:"spec"`
}

type checkRepositoryInput struct {
	Name        string    			`json:"name"`
	Source		string				`json:"source"`
	Path		string				`json:"path"`
	URL			string				`json:"url"`
	KeysetId	*int				`json:"keyset"`
	IsPrivate 	bool				`json:"is_private"`
	Config 		map[string]*string	 `json:"config"`
}

type checkRepositoryResponse struct {
	Accessible 		bool    	`json:"accessible"`
	Directories		[]string 	`json:"directories"`
	IsMultiChart	bool		`json:"is_multi_chart"`
	Error			*string		`json:"error"`
	Contents		[]checkRepositoryContents	`json:"contents"`
}

type checkRepositoryContents struct {
	Name          string    	`json:"name"`
	URL			  string		`json:"url"`
	HtmlUrl		  string		`json:"html_url"`
	DownloadURL	  string		`json:"download_url"`
	Sha			  string		`json:"sha"`
	Links		  map[string]string	`json:"_links"`
	GitURL		string			`json:"git_url"`
	Path		string			`json:"path"`
	Type		string			`json:"type"`
	Size		int				`json:"size"`
}

var repositoryCmd = &cobra.Command{
	Use:   "repositories",
	Aliases: []string{"repos", "repo"},
	Short: "manage chart repositories",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("NKS CLI Version - v0.0.1 (alpha)")
	//},
}

var listRepositoryCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"l", "li", "lis"},
	Short: "list custom repositories",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		repos := listRepositories()
		printRepositories(repos)
	},
}

func listRepositories() []createRepoResponse {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos", CurrentConfig.OrgID)
	res, err := httpRequest("GET", url)
	check(err)

	data := []createRepoResponse{}

	err = json.Unmarshal(res, &data)
	check(err)

	return data
}

var createRepositoryCmd = &cobra.Command{
	Use:   "create",
	Aliases: []string{"cr", "crea"},
	Short: "create a custom chart repo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.HasPrefix(flagRepositoryURL, "github.com") {
			flagRepositoryURL = fmt.Sprintf("https://%s", flagRepositoryURL)
		}
		i := checkRepositoryInput{
			Name: flagRepositoryName,
			URL: flagRepositoryURL,
			Source: flagCreateRepositorySourceType,
		}

		_, err := checkRepository(i)
		check(err)

		input := createRepoInput{i}

		n, err := createRepository(input)
		check(err)
		

		printRepositories(*n)
	},
}

var flagRepositoryName string
var flagRepositorySource string
var flagRepositoryURL	string
var flagRepositoryID	int


func checkRepository(i checkRepositoryInput) (*checkRepositoryResponse, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos/check", CurrentConfig.OrgID)
	b, err := json.Marshal(i)
	check(err)

	//fmt.Printf("body1: %s\n\n", string(b))
	res, err := httpRequestPost("POST", url, b)
	check(err)
	//fmt.Printf("response1: %s\n\n", string(res))

	data := checkRepositoryResponse{}

	err = json.Unmarshal(res, &data)
	check(err)

	return &data, err
}

func createRepository(i createRepoInput) (*createRepoResponseS, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos", CurrentConfig.OrgID)
	b, err := json.Marshal(i)
	check(err)

	//fmt.Printf("body2: %s\n\n", string(b))
	res, err := httpRequestPost("POST", url, b)
	check(err)
	//fmt.Printf("response2: %s\n\n", string(res))

	data := createRepoResponseS{}

	err = json.Unmarshal(res, &data)
	check(err)

	return &data, err
}

var deleteRepositoryCmd = &cobra.Command{
	Use:   "delete",
	Aliases: []string{"rm", "del", "remove"},
	Short: "delete a custom chart repo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		deleteRepository(flagRepositoryID)
		printRepositories(listRepositories())
	},
}

func deleteRepository(repoId int) (error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos/%d", CurrentConfig.OrgID, repoId)
	
	_, err := httpRequest("DELETE", url)
	check(err)

	return err
}

func printRepositories(rs []createRepoResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tSOURCE\tURL\t# CHARTS\t\n")
	for _, c := range rs {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", c.Name, c.ID, c.SourceDisplay, c.URL, len(c.ChartIndex))
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(repositoryCmd)

	repositoryCmd.AddCommand(createRepositoryCmd)
	repositoryCmd.AddCommand(listRepositoryCmd)
	repositoryCmd.AddCommand(deleteRepositoryCmd)
	repositoryCmd.Flags().StringVarP(&flagCreateRepositorySourceType, "type", "t", "github", "repository type")

	createRepositoryCmd.Flags().StringVarP(&flagRepositoryName, "name", "n", "", "name of repository")
	createRepositoryCmd.Flags().StringVarP(&flagRepositorySource, "source", "s", "github", "repository source (default: github)")
	createRepositoryCmd.Flags().StringVarP(&flagRepositoryURL, "url", "u", "", "url of repository")

	deleteRepositoryCmd.Flags().IntVarP(&flagRepositoryID, "id", "i", 0, "id of repository")
	
	e := createRepositoryCmd.MarkFlagRequired("name")
	check(e)
	e = createRepositoryCmd.MarkFlagRequired("url")
	check(e)
	e = deleteRepositoryCmd.MarkFlagRequired("id")
	check(e)
}
