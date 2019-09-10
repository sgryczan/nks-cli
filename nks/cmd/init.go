package cmd

import (
	"fmt"

	nks "github.com/NetApp/nks-sdk-go/nks"
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

var configBootStrap bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !configBootStrap {
			newConfig()
		}
	},
}

func init() {
	configCmd.AddCommand(initCmd)
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
	o, err := getDefaultOrg()
	if err != nil {
		fmt.Errorf("Could not get default organization")
	} else {
		fmt.Println("Setting default org..")
		setConfigDefaultOrg(o)
	}
	viper.WriteConfigAs(filename)

	return nil
}

func setConfigDefaultOrg(o nks.Organization) {
	viper.Set("org_id", o.ID)
	viper.WriteConfig()
}

func bootstrapConfigFile() {
	configBootStrap = true
	newConfig()
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
