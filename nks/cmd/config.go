/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

type config struct {
	OrgID            int    `mapstructure:"org_id"`
	Provider         string `mapstructure:"provider"`
	ProviderKeySetID int    `mapstructure:"provider_keyset_id"`
	ApiToken         string `mapstructure:"api_token"`
	ApiURL           string `mapstructure:"api_url"`
	ClusterId        int    `mapstructure:"cluster_id"`
	SSHKeySetId      int    `mapstructure:"ssh_keyset_id"`
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "nks cli configuration",
	Long:  `Various commands for configuring nks cli`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("config called")
	//},
}

var listConfigCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "list current configuration",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range viper.AllSettings() {
			fmt.Printf("%s: %v\n", k, v)
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set configuration setting",
	Long:  "",
	/* Run: func(cmd *cobra.Command, args []string) {
		for k, v := range viper.AllSettings() {
			fmt.Printf("%s: %v\n", k, v)
		}
	}, */
}

var configSetTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "set api token",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("api_token", args[0])
		viper.WriteConfig()
	},
}

var configSetURLCmd = &cobra.Command{
	Use:   "url",
	Short: "set api url",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("api_url", args[0])
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(listConfigCmd)
	configCmd.AddCommand(configSetCmd)
	configSetCmd.AddCommand(configSetTokenCmd)
	configSetCmd.AddCommand(configSetURLCmd)
}
