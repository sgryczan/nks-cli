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
	"bufio"
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		newConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createConfigFile(filename string, apitoken string) error {
	s := fmt.Sprintf("api_token: %s\n", apitoken)
	bs := []byte(s)

	f, err := os.Create(filename)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(bs)
	if err != nil {
		panic(err)
	}

	f.Sync()

	return nil
}

func readApiToken() string {
	fmt.Printf("Enter your NKS API Token:")
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')

	// convert CRLF to LF
	token = strings.Replace(token, "\n", "", -1)

	return token
}

func newConfig() error {
	home, _ := homedir.Dir()

	fmt.Println("Creating config file...")
	token := readApiToken()

	createConfigFile(fmt.Sprintf("%s/.nks.yaml", home), token)

	// Search config in home directory with name ".nks" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".nks.yaml")

	return nil
}
