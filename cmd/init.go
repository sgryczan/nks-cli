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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFields = []string{
	"api_token",
	"api_url",
	"org_id",
	"cluster_id",
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new configuration",
	Run: func(cmd *cobra.Command, args []string) {
		newConfig()
	},
}

func init() {
	configCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createConfigFile(filename string, token string) error {
	for _, s := range configFields {
		var v string
		switch s {
		case "api_token":
			v = token
		case "api_url":
			v = "https://api.nks.netapp.io"
		default:

		}
		viper.Set(s, v)
	}
	fmt.Println("Setting default org..")
	setConfigDefaultOrg(getDefaultOrg())
	viper.WriteConfigAs(filename)

	return nil
}

func setConfigDefaultOrg(o organization) {
	viper.Set("org_id", o.ID)
	viper.WriteConfig()
}

func bootstrapConfigFile() {

}

func newConfig() error {
	home, _ := homedir.Dir()

	fmt.Println("Creating config file...")
	token := readApiToken()

	createConfigFile(fmt.Sprintf("%s/.nks.yaml", home), token)

	// Search config in home directory with name ".nks" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".nks")

	return nil
}
