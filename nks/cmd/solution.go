package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
	models "gitlab.com/sgryczan/nks-cli/nks/models"
)

var solutionsCmd = &cobra.Command{
	Use:     "solutions",
	Aliases: []string{"sol", "solu", "solution", "so", "chart", "charts"},
	Short:   "mnanage solutions",
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

var listSolutionsCmd = &cobra.Command{
	Use:     "list-installed",
	Aliases: []string{"li"},
	Short:   "list solutions",
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		clusterId := vpr.GetInt("cluster_id")
		orgId := vpr.GetInt("org_id")
		if flagClusterId != 0 {
			clusterId = flagClusterId
		}

		if clusterId == 0 {
			fmt.Printf("No default cluster set. Set one, or specify a cluster with --cluster-id.\n")
			os.Exit(1)
		}

		ss, err := listSolutions(orgId, clusterId)
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			ss = &[]nks.Solution{}
		}
		fmt.Printf("(cluster id: %d)\n\n", clusterId)
		printSolutions(*ss)
	},
}

var listSolutionTemplatesCmd = &cobra.Command{
	Use:     "list-templates",
	Aliases: []string{"lt"},
	Short:   "list solution templates",
	Run: func(cmd *cobra.Command, args []string) {
		s := models.ListSolutionTemplates()

		models.PrintSolutionTemplates(&s)
	},
}

var deploySolutionFromTemplateCmd = &cobra.Command{
	Use:     "deploy-template",
	Aliases: []string{"dt"},
	Short:   "deploy solution (jenkins)",
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		name := "jenkins"
		fmt.Printf("creating solution %s...\n", name)
		cid := vpr.GetInt("cluster_id")

		if flagClusterId != 0 {
			cid = flagClusterId
		}

		s, err := createSolutionFromTemplate(name, vpr.GetInt("org_id"), cid)

		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
		}
		ss := []nks.Solution{*s}
		printSolutions(ss)
	},
}

var deploySolutionFromRepositoryCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"new", "dep"},
	Short:   "deploy an imported chart",
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		repoName := flagSolutionRepoName

		if FlagDebug {
			fmt.Printf("creating solution %s...\n", repoName)
		}

		cid := vpr.GetInt("cluster_id")

		if flagClusterId != 0 {
			cid = flagClusterId
		}

		repo, err := GetRepositoryByName(repoName, FlagDebug)
		if err != nil {
			fmt.Printf("We had an error trying to retrieve Repository: %s\n", repoName)
			os.Exit(1)
		}

		if FlagDebug {
			fmt.Printf("Converting Repositiory:\n %+v\n", repo)
		}
		s, err := createSolutionFromRepository(repo, repoName, vpr.GetInt("org_id"), cid)

		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
		}
		ss := []nks.Solution{*s}
		printSolutions(ss)
	},
}

var deleteSolutionsCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "delete solution",
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		clusterId := vpr.GetInt("cluster_id")
		orgId := vpr.GetInt("org_id")
		if flagClusterId != 0 {
			clusterId = flagClusterId
		}

		if clusterId == 0 {
			fmt.Printf("No default cluster set. Set one, or specify a cluster with --cluster-id.\n")
			os.Exit(1)
		}

		fmt.Printf("deleting solution %d...\n", flagSolutionId)

		err := deleteSolution(orgId, clusterId, flagSolutionId)
		check(err)

		ss, err := listSolutions(orgId, clusterId)
		printSolutions(*ss)
	},
}

func listSolutions(orgId, clusterId int) (*[]nks.Solution, error) {

	s, err := SDKClient.GetSolutions(orgId, clusterId)
	check(err)

	return &s, err
}

func printSolutions(s []nks.Solution) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tSOLUTION\tSTATE\t\n")
	for _, s := range s {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", s.Name, s.ID, s.Solution, s.State)
	}
	w.Flush()
}

func createSolutionFromTemplate(s string, orgId, clusterId int) (*nks.Solution, error) {

	template, err := models.GetTemplateAsJson(s)
	check(err)

	//fmt.Printf("Solution body: %s\n", template)
	sol, err := SDKClient.AddSolutionFromJSON(orgId, clusterId, template)
	return sol, err
}

func createSolutionFromRepository(r models.Repository, releaseName string, orgId, clusterId int) (*nks.Solution, error) {

	template := models.RepositoryToTemplate(r, releaseName)
	b, err := json.Marshal(template)
	if err != nil {
		fmt.Printf("error while attempting to convert template: \n\t:%v", err)
	}

	if FlagDebug {
		fmt.Printf("Solution body: %s\n", string(b))
	}
	sol, err := SDKClient.AddSolutionFromJSON(orgId, clusterId, string(b))
	return sol, err
}

func deleteSolution(orgId, clusterId, solutionId int) error {

	err := SDKClient.DeleteSolution(orgId, clusterId, solutionId)
	return err
}

func init() {
	rootCmd.AddCommand(solutionsCmd)
	solutionsCmd.AddCommand(listSolutionsCmd)
	solutionsCmd.AddCommand(listSolutionTemplatesCmd)
	solutionsCmd.AddCommand(deploySolutionFromTemplateCmd)
	solutionsCmd.AddCommand(deploySolutionFromRepositoryCmd)
	solutionsCmd.AddCommand(deleteSolutionsCmd)

	deploySolutionFromRepositoryCmd.Flags().StringVarP(&flagSolutionRepoName, "name", "n", "demo", "Name of target repository")
	deploySolutionFromTemplateCmd.Flags().IntVarP(&flagClusterId, "cluster-id", "c", vpr.GetInt("cluster_id"), "ID of target cluster")
	deploySolutionFromTemplateCmd.Flags().StringVarP(&flagSolutionName, "name", "n", "jenkins", "Name of solution template")

	listSolutionsCmd.Flags().IntVarP(&flagClusterId, "cluster-id", "c", vpr.GetInt("cluster_id"), "ID of target cluster")

	deleteSolutionsCmd.Flags().IntVarP(&flagClusterId, "cluster-id", "c", vpr.GetInt("cluster_id"), "ID of target cluster")
	deleteSolutionsCmd.Flags().IntVarP(&flagSolutionId, "solution-id", "s", 0, "ID of solution")
	e := deleteSolutionsCmd.MarkFlagRequired("solution-id")
	check(e)
}
