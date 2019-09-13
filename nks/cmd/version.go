package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display CLI version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NKS CLI Version - v0.0.1 (alpha)")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
