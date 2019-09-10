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

// workspacesCmd represents the workspaces command
var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	/* Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspaces called")
	}, */
}

func init() {
}
