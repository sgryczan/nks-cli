package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
	homedir "github.com/mitchellh/go-homedir"
	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/cobra"
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

type config struct {
	OrgID            int    `mapstructure:"org_id"`
	Provider         string `mapstructure:"provider"`
	ProviderKeySetID int    `mapstructure:"provider_keyset_id"`
	ApiToken         string `mapstructure:"api_token"`
	ApiURL           string `mapstructure:"api_url"`
	ClusterId        int    `mapstructure:"cluster_id"`
	SSHKeySetId      int    `mapstructure:"ssh_keyset_id"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Aliases: []string{"conf"},
	Short: "nks cli configuration",
	Long:  `Various commands for configuring nks cli`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("config called")
	//},
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !configBootStrap {
			newConfig()
		}
	},
}

var listConfigCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "list current configuration",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current configuration:\n\n")
		t := CurrentConfig
		s := reflect.ValueOf(t).Elem()
		typeOfT := s.Type()

		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Printf("%s %s = %v\n",
				typeOfT.Field(i).Name, f.Type(), f.Interface())
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
		syncConfig()
	},
}

func configSet(key string, value string) {
	viper.Set(key, value)
	syncConfig()
}

var configSetClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "set default cluser",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		setCluster(flagClusterId)
	},
}

func setCluster(clusterId int) {
	viper.Set("cluster_id", clusterId)
	syncConfig()
	CurrentConfig.ClusterId = clusterId
}

var configSetURLCmd = &cobra.Command{
	Use:   "url",
	Short: "set api url",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("api_url", args[0])
		syncConfig()
	},
}

func syncConfig() {
	viper.WriteConfig()
	err := viper.Unmarshal(CurrentConfig)
	check(err)
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

	return nil
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
	configureDefaultOrganization()
	viper.WriteConfigAs(filename)

	fmt.Println("Setting Provider Key...")
	setDefaultProviderKey(viper.GetString("provider"))

	fmt.Println("Setting SSH Key...")
	setSSHKey(viper.GetString("provider"))

	return nil
}

func configureDefaultOrganization() {
	o, err := getDefaultOrg()
	if err != nil {
		fmt.Println("Could not get default organization")
	} else {
		fmt.Println("Setting default org..")
		setConfigDefaultOrg(o)
	}
}

func setConfigDefaultOrg(o nks.Organization) {
	viper.Set("org_id", o.ID)
	syncConfig()
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
	syncConfig()
}

func setDefaultSSHKey() {
	kss, err := getKeySets()
	check(err)
	v := []nks.Keyset{}

	for _, ks := range *kss {
		if ks.Category == "user_ssh" {
			v = append(v, ks)
		}
	}

	if len(v) == 0 {
		fmt.Println("Error configuring default SSH keyset. No user ssh keysets found!")
	}

	viper.Set("ssh_keyset_id", v[0].ID)
	syncConfig()
}

func setSSHKey(s string) {
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
	syncConfig()
}

func bootstrapConfigFile() {
	configBootStrap = true
	newConfig()
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(listConfigCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(initConfigCmd)
	configSetCmd.AddCommand(configSetTokenCmd)
	configSetCmd.AddCommand(configSetURLCmd)
	configSetCmd.AddCommand(configSetClusterCmd)

	configSetClusterCmd.Flags().IntVarP(&flagClusterId, "id", "i", CurrentConfig.ClusterId, "ID of target cluster")
	e := configSetClusterCmd.MarkFlagRequired("id")
	check(e)
}
