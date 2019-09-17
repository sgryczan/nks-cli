package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	models "gitlab.com/sgryczan/nks-cli/nks/models"

	"github.com/spf13/cobra"
)

var flagCreateRepositorySourceType string

var repositoryCmd = &cobra.Command{
	Use:     "repositories",
	Aliases: []string{"repos", "repo"},
	Short:   "manage chart repositories",
	Long:    ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("NKS CLI Version - v0.0.1 (alpha)")
	//},
}

var listRepositoryCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li", "lis"},
	Short:   "list custom repositories",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		repos := listRepositories()
		printRepositories(repos)
	},
}

func listRepositories() []models.Repository {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos", CurrentConfig.OrgID)
	res, err := httpRequest("GET", url)
	check(err)

	data := []models.Repository{}

	err = json.Unmarshal(res, &data)
	check(err)

	return data
}

func GetRepositoryByName(name string, debug bool) (models.Repository, error) {
	var err error
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos", CurrentConfig.OrgID)
	res, err := httpRequest("GET", url)
	check(err)

	repositories := []models.Repository{}

	err = json.Unmarshal(res, &repositories)
	check(err)

	for _, repo := range repositories {

		if debug {
			fmt.Printf("GetRepositoryByName() - Checking repository %s\n", name)
		}
		if repo.Name == name {
			if debug {
				fmt.Printf("GetRepositoryByName() - Matched repository %s\n", name)
			}
			return repo, err
		}
	}
	if FlagDebug {
		fmt.Printf("GetRepositoryByName() - Failed to match repository %s\n", name)
	}
	r := models.Repository{}
	err = fmt.Errorf("Failed to match repository - %s", name)
	return r, err
}

var createRepositoryCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr", "crea"},
	Short:   "create a custom chart repo",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.HasPrefix(flagRepositoryURL, "github.com") {
			flagRepositoryURL = fmt.Sprintf("https://%s", flagRepositoryURL)
		}
		i := models.CheckRepositoryInput{
			Name:   flagRepositoryName,
			URL:    flagRepositoryURL,
			Source: flagCreateRepositorySourceType,
		}

		_, err := checkRepository(i)
		check(err)

		input := models.CreateRepoInput{i}

		n, err := createRepository(input)
		check(err)

		printRepositories(*n)
	},
}

var flagRepositoryName string
var flagRepositorySource string
var flagRepositoryURL string
var flagRepositoryID int

func checkRepository(i models.CheckRepositoryInput) (*models.CheckRepositoryResponse, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos/check", CurrentConfig.OrgID)
	b, err := json.Marshal(i)
	check(err)

	//fmt.Printf("body1: %s\n\n", string(b))
	res, err := httpRequestPost("POST", url, b)
	check(err)
	//fmt.Printf("response1: %s\n\n", string(res))

	data := models.CheckRepositoryResponse{}

	err = json.Unmarshal(res, &data)
	check(err)

	return &data, err
}

func createRepository(i models.CreateRepoInput) (*models.RepositoryS, error) {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos", CurrentConfig.OrgID)
	b, err := json.Marshal(i)
	check(err)

	//fmt.Printf("body2: %s\n\n", string(b))
	res, err := httpRequestPost("POST", url, b)
	check(err)
	//fmt.Printf("response2: %s\n\n", string(res))

	data := models.RepositoryS{}

	err = json.Unmarshal(res, &data)
	check(err)

	return &data, err
}

var deleteRepositoryCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del", "remove"},
	Short:   "delete a custom chart repo",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		deleteRepository(flagRepositoryID)
		printRepositories(listRepositories())
	},
}

func deleteRepository(repoId int) error {
	url := fmt.Sprintf("https://api.nks.netapp.io/orgs/%d/chart-repos/%d", CurrentConfig.OrgID, repoId)

	_, err := httpRequest("DELETE", url)
	check(err)

	return err
}

func printRepositories(rs []models.Repository) {
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
