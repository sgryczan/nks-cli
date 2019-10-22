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

		clusterID := vpr.GetInt("cluster_id")
		orgID := vpr.GetInt("org_id")
		if flagclusterID != 0 {
			clusterID = flagclusterID
		}

		if clusterID == 0 {
			fmt.Printf("No default cluster set. Set one, or specify a cluster with --cluster-id.\n")
			os.Exit(1)
		}

		ss, err := listSolutions(orgID, clusterID)
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
			ss = &[]nks.Solution{}
		}
		fmt.Printf("(cluster id: %d)\n\n", clusterID)
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

		if flagclusterID != 0 {
			cid = flagclusterID
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

		if flagDebug {
			fmt.Printf("creating solution %s...\n", repoName)
		}

		cid := vpr.GetInt("cluster_id")

		if flagclusterID != 0 {
			cid = flagclusterID
		}

		repo, err := GetRepositoryByName(repoName, flagDebug)
		if err != nil {
			fmt.Printf("We had an error trying to retrieve Repository: %s\n", repoName)
			os.Exit(1)
		}

		if flagDebug {
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

		clusterID := vpr.GetInt("cluster_id")
		orgID := vpr.GetInt("org_id")
		if flagclusterID != 0 {
			clusterID = flagclusterID
		}

		if clusterID == 0 {
			fmt.Printf("No default cluster set. Set one, or specify a cluster with --cluster-id.\n")
			os.Exit(1)
		}

		fmt.Printf("deleting solution %d...\n", flagsolutionID)

		err := deleteSolution(orgID, clusterID, flagsolutionID)
		check(err)

		ss, err := listSolutions(orgID, clusterID)
		printSolutions(*ss)
	},
}

func listSolutions(orgID, clusterID int) (*[]nks.Solution, error) {

	s, err := SDKClient.GetSolutions(orgID, clusterID)
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

func createSolutionFromTemplate(s string, orgID, clusterID int) (*nks.Solution, error) {

	template, err := models.GetTemplateAsJson(s)
	check(err)

	//fmt.Printf("Solution body: %s\n", template)
	sol, err := SDKClient.AddSolutionFromJSON(orgID, clusterID, template)
	return sol, err
}

func createSolutionFromRepository(r models.Repository, releaseName string, orgID, clusterID int) (*nks.Solution, error) {

	template := models.RepositoryToTemplate(r, releaseName)
	b, err := json.Marshal(template)
	if err != nil {
		fmt.Printf("error while attempting to convert template: \n\t:%v", err)
	}

	if flagDebug {
		fmt.Printf("Solution body: %s\n", string(b))
	}
	sol, err := SDKClient.AddSolutionFromJSON(orgID, clusterID, string(b))
	return sol, err
}

func deleteSolution(orgID, clusterID, solutionID int) error {

	err := SDKClient.DeleteSolution(orgID, clusterID, solutionID)
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
	deploySolutionFromTemplateCmd.Flags().IntVarP(&flagclusterID, "cluster-id", "c", vpr.GetInt("cluster_id"), "ID of target cluster")
	deploySolutionFromTemplateCmd.Flags().StringVarP(&flagSolutionName, "name", "n", "jenkins", "Name of solution template")

	listSolutionsCmd.Flags().IntVarP(&flagclusterID, "cluster-id", "c", vpr.GetInt("cluster_id"), "ID of target cluster")

	deleteSolutionsCmd.Flags().IntVarP(&flagclusterID, "cluster-id", "c", vpr.GetInt("cluster_id"), "ID of target cluster")
	deleteSolutionsCmd.Flags().IntVarP(&flagsolutionID, "solution-id", "s", 0, "ID of solution")
	e := deleteSolutionsCmd.MarkFlagRequired("solution-id")
	check(e)
}
