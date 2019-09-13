package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	models "gitlab.com/sgryczan/nks-cli/nks/models"
	nks "github.com/NetApp/nks-sdk-go/nks"
)

var solutionsCmd = &cobra.Command{
	Use:   "solutions",
	Aliases: []string{"sol", "solu", "solution", "so", "chart", "charts"},
	Short: "mnanage solutions",
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

var listSolutionsCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"l", "li", "lis"},
	Short: "list solutions",
	Run: func(cmd *cobra.Command, args []string) {
		clusterId := CurrentConfig.ClusterId
		orgId := CurrentConfig.OrgID
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

var createSolutionsCmd = &cobra.Command{
	Use:   "deploy",
	Aliases: []string{"new", "dep"},
	Short: "deploy solution (jenkins)",
	Run: func(cmd *cobra.Command, args []string) {
		name := "jenkins"
		fmt.Printf("creating solution %s...\n", name)
		cid := CurrentConfig.ClusterId

		if flagClusterId != 0 {
			cid = flagClusterId
		}	

		s, err := createSolution(name, CurrentConfig.OrgID, cid)
		
		if err != nil {
			fmt.Printf("There was an error retrieving items:\n\t%s\n\n", err)
		}
		ss := []nks.Solution{*s,}
		printSolutions(ss)
	},
}

var deleteSolutionsCmd = &cobra.Command{
	Use:   "delete",
	Aliases: []string{"rm", "del"},
	Short: "delete solution",
	Run: func(cmd *cobra.Command, args []string) {
		clusterId := CurrentConfig.ClusterId
		orgId := CurrentConfig.OrgID
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

	c := newClient()
	s, err := c.GetSolutions(orgId, clusterId)
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

func createSolution(s string, orgId, clusterId int) (*nks.Solution, error) {

	template, err := models.GetTemplateAsJson(s)
	check(err)

	c := newClient()
	//fmt.Printf("Solution body: %s\n", template)
	sol, err := c.AddSolutionFromJSON(orgId, clusterId, template)
	return sol, err
}

func deleteSolution(orgId, clusterId, solutionId int) (error) {

	c := newClient()
	err := c.DeleteSolution(orgId, clusterId, solutionId)
	return err
}

func init() {
	rootCmd.AddCommand(solutionsCmd)
	solutionsCmd.AddCommand(listSolutionsCmd)
	solutionsCmd.AddCommand(createSolutionsCmd)
	solutionsCmd.AddCommand(deleteSolutionsCmd)

	createSolutionsCmd.Flags().IntVarP(&flagClusterId, "cluster-id", "c", CurrentConfig.ClusterId, "ID of target cluster")
	listSolutionsCmd.Flags().IntVarP(&flagClusterId, "cluster-id", "c", CurrentConfig.ClusterId, "ID of target cluster")

	deleteSolutionsCmd.Flags().IntVarP(&flagClusterId, "cluster-id", "c", CurrentConfig.ClusterId, "ID of target cluster")
	deleteSolutionsCmd.Flags().IntVarP(&flagSolutionId, "solution-id", "s", 0, "ID of solution")
	e := deleteSolutionsCmd.MarkFlagRequired("solution-id")
	check(e)
}
