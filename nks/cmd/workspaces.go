package cmd

import (
	"github.com/spf13/cobra"
)

type Workspace struct {
	ID          int    `json:"pk,Number"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Org         string `json:"org"`
	Default     string `json:"is_default"`
	Pinned      string `json:"is_pinned"`
	CreatedTime string `json:"created"`
}

var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "manage workspaces",
	Long: ``,
	/* Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspaces called")
	}, */
}

func init() {
}
