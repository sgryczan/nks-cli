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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type organization struct {
	ID                 json.Number `json:"pk,Number"`
	Name               string      `json:"name"`
	Slug               string      `json:"slug"`
	Logo               string      `json:"logo"`
	EnableExperimental string      `json:"enable_experimental_feature,Bool"`
	CreatedTime        string      `json:"created"`
	UpdatedTime        string      `json:"updated"`
}

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("organization called")
	//},
}

func printOrgs(o []organization) {
	fmt.Printf("NAME\t\tID\t\n")
	for _, org := range o {
		fmt.Printf("%v\t%v\t\n", org.Name, org.ID)
	}
}

func getOrgs() *[]organization {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.nks.netapp.io/orgs", nil)
	req.Header.Add("Authorization", "Bearer "+viper.GetString("api_token"))

	resp, err := client.Do(req)
	check(err)

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	data := []organization{}

	err = json.Unmarshal(body, &data)
	check(err)

	return &data
}

func getDefaultOrg() organization {
	o := getOrgs()

	// First Org returned is the default
	return (*o)[0]
}

// organizationCmd represents the organization command
var organizationGetCmd = &cobra.Command{
	Use:   "orgs",
	Short: "Get Organizations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		data := getOrgs()
		printOrgs(*data)
	},
}

func init() {
	getCmd.AddCommand(organizationGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// organizationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// organizationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
