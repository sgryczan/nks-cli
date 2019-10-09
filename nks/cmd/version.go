package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display CLI version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion(version)
	},
}

func printVersion(v string) {
	fmt.Printf("NKS CLI Version - v%s (alpha)\n", v)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
