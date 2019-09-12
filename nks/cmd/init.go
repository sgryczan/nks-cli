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
	"provider",
	"provider_keyset_id",
	"ssh_keyset_id",
}

var CurrentConfig = &config{}

var configBootStrap bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !configBootStrap {
			newConfig()
		}
	},
}

func createConfigFile(filename string, token string) error {
	for _, s := range configFields {
		var v interface{}
		switch s {
		case "api_token":
			v = token
		case "api_url":
			v = "https://api.nks.netapp.io"
		case "provider":
			v = "gce"
		case "cluster_id":
			v = 0
		default:

		}
		viper.Set(s, v)
	}
	o, err := getDefaultOrg()
	if err != nil {
		fmt.Println("Could not get default organization")
	} else {
		fmt.Println("Setting default org..")
		setConfigDefaultOrg(o)
	}
	viper.WriteConfigAs(filename)

	fmt.Println("Setting Provider Key...")
	setDefaultProviderKey(viper.GetString("provider"))

	fmt.Println("Setting SSH Key...")
	setDefaultSSHKey(viper.GetString("provider"))

	return nil
}

func setConfigDefaultOrg(o nks.Organization) {
	viper.Set("org_id", o.ID)
	viper.WriteConfig()
}

func setDefaultProviderKey(p string) {
	kss, err := getKeySets()
	check(err)
	v := []nks.Keyset{}

	for _, ks := range *kss {
		if ks.Entity == p {
			v = append(v, ks)
		}
	}

	if len(v) == 0 {
		fmt.Printf("no keysets found for provider %s!\n", p)
	}

	viper.Set("provider_keyset_id", v[0].ID)
	viper.WriteConfig()
}

func setDefaultSSHKey(s string) {
	kss, err := getKeySets()
	check(err)
	v := []nks.Keyset{}

	for _, ks := range *kss {
		if ks.Category == "user_ssh" {
			v = append(v, ks)
		}
	}

	if len(v) == 0 {
		fmt.Println("No user ssh keysets found")
	}

	viper.Set("ssh_keyset_id", v[0].ID)
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

	err := viper.Unmarshal(CurrentConfig)
	check(err)
	for k, v := range viper.AllSettings() {
		fmt.Printf("%s: %v\n", k, v)
	}
	fmt.Printf("Config: %+v", CurrentConfig)

	return nil
}

func init() {
	configCmd.AddCommand(initCmd)
}
