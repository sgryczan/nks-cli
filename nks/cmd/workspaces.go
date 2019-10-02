package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
)

var workspacesCmd = &cobra.Command{
	Use:     "workspaces",
	Aliases: []string{"ws"},
	Short:   "manage workspaces",
	Long:    ``,
	/* Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspaces called")
	}, */
}

var listWorkspacesCmd = &cobra.Command{
	Use:   "list",
	Short: "list workspaces in organization",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		checkDefaultOrg()

		o := vpr.GetInt("org_id")
		if flagOrganizationId != 0 {
			o = flagOrganizationId
		}
		ws, err := listWorkspaces(o)
		if err != nil {
			fmt.Printf("There was an error retrieving workspaces: %v\n", err)
			os.Exit(1)
		}

		printWorkspaces(ws)
	},
}

func listWorkspaces(orgId int) ([]nks.Workspace, error) {

	ws, err := SDKClient.GetWorkspaces(orgId)
	return ws, err
}

func printWorkspaces(wss []nks.Workspace) {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fmt.Fprintf(w, "NAME\tID\tORG\tCLUSTERS\tFEDERATIONS\tTEAM WORKSPACES\t\n")
	for _, ws := range wss {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t\n", ws.Name, ws.ID, ws.Org, len(ws.Clusters), len(ws.Federations), len(ws.TeamWorkspaces))
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(workspacesCmd)
	workspacesCmd.AddCommand(listWorkspacesCmd)
	listWorkspacesCmd.Flags().IntVarP(&flagOrganizationId, "organization-id", "i", 0, "Organization ID")
}
