package cmd

import (
	"fmt"
	"os"

	nks "github.com/NetApp/nks-sdk-go/nks"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	vpr "github.com/spf13/viper"
	ext "gitlab.com/sgryczan/nks-cli/nks/extensions"
)

var configFields = []string{
	"api_token",
	"api_url",
	"org_id",
	"workspace_id",
	"cluster_id",
	"hci_keyset",
	"aws_keyset",
	"gce_keyset",
	"gke_keyset",
	"azr_keyset",
	"eks_keyset",
	"ssh_keyset",
	"provider_keyset_id",
	"provider",
}

var configBootStrap bool

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "nks cli configuration",
	Long:    `Various commands for configuring nks cli`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("config called")
	//},
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "create a new configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !configBootStrap {
			newConfigFile()
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

		for _, k := range vpr.AllKeys() {
			v := vpr.Get(k)
			if v != nil {
				fmt.Printf("%s = %v\n",
					k, v)
			}

		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set configuration setting",
	Long:  "",
	/* Run: func(cmd *cobra.Command, args []string) {
		for k, v := range vpr.AllSettings() {
			fmt.Printf("%s: %v\n", k, v)
		}
	}, */
}

var configSetTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "set api token",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		vpr.Set("api_token", args[0])
	},
}

func configSet(key string, value string) {
	vpr.Set(key, value)
}

var configSetClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "set default cluser",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		setClusterAsCurrent(flagclusterID)
	},
}

var configSetOrganizationCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"org"},
	Short:   "set default organization",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		setOrgID(flagOrganizationID)
		setclusterID(0)
		clearCurrentKeySets()
		syncConfigFile()
	},
}

var configSetWorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "set default workspace",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		setWorkSpaceID(flagWorkspaceID)
		syncConfigFile()
	},
}

func setClusterAsCurrent(clusterID int) {
	if flagDebug {
		fmt.Printf("Debug - setClusterAsCurrent(%d)\n", clusterID)
	}
	setClusterKubeConfig(clusterID)
	setclusterID(clusterID)
}

func setclusterID(clusterID int) {
	if flagDebug {
		fmt.Printf("Debug - setclusterID(%d)\n", clusterID)
	}
	vpr.Set("cluster_id", clusterID)
	vpr.WriteConfig()
}

func setOrgID(orgID int) {
	if flagDebug {
		fmt.Printf("Debug - setOrgID(%d)\n", orgID)
	}
	vpr.Set("org_id", orgID)
}

func setWorkSpaceID(workspaceID int) {
	vpr.Set("workspace_id", workspaceID)
}

var configSetURLCmd = &cobra.Command{
	Use:   "url",
	Short: "set api url",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		vpr.Set("api_url", args[0])

	},
}

func clearCurrentKeySets() {
	if flagDebug {
		fmt.Println("Debug - clearCurrentKeySets()")
	}

	keysets := []string{
		"hci_keyset",
		"aws_keyset",
		"gce_keyset",
		"gke_keyset",
		"azr_keyset",
		"eks_keyset",
		"ssh_keyset",
	}

	for _, k := range keysets {
		vpr.Set(k, 0)
	}
}

func syncConfigFile() {
	if flagDebug {
		fmt.Println("Debug - syncConfigFile()")
	}
	if err := vpr.ReadInConfig(); err != nil {
		if flagDebug {
			fmt.Println("syncConfigFile() - No config file detected. Creating to allow sync.")
		}
		newConfigFile()
	}
	vpr.WriteConfig()
}

func newConfigFile() error {
	if flagDebug {
		fmt.Println("Debug - newConfigFile()")
	}
	home, _ := homedir.Dir()
	var token string

	fmt.Println("Creating config file...")
	if t := vpr.GetString("api_token"); t != "" {
		token = t
	} else {
		token = readAPIToken()
		vpr.Set("api_token", token)
	}

	createConfigFile(fmt.Sprintf("%s/.nks.yaml", home), token)

	if flagSetDefaults {
		if vpr.GetInt("org_id") == 0 {
			configureDefaultOrganization()
		}

		if vpr.GetInt("ssh_keyset") == 0 {
			fmt.Println("Default SSH keys not set. Configuring...")
			setDefaultProviderKey("user_ssh")
		}

		//if vpr.GetString("provider") == "" {
		//	configSet("provider", "gce")
		//}

		//if vpr.GetInt("provider_keyset_id") == 0 {
		//	setDefaultProviderKey(vpr.GetString("provider"))
		//}
	}

	// Search config in home directory with name ".nks" (without extension).
	vpr.AddConfigPath(home)
	vpr.SetConfigName(".nks")

	return nil
}

func createConfigFile(filename string, token string) error {
	if flagDebug {
		fmt.Println("Debug - createConfigFile()")
	}
	for _, s := range configFields {
		var v interface{}
		switch s {
		case "api_token":
			v = token
		case "api_url":
			v = "https://api.nks.netapp.io"
		default:
		}
		vpr.Set(s, v)
	}

	vpr.WriteConfigAs(filename)

	return nil
}

func configureDefaultOrganization() {

	if flagDebug {
		fmt.Println("Debug - configureDefaultOrganization()")
	}

	if SDKClient.Token == "" {
		initClient()
	}
	profile, err := SDKClient.GetUserProfile()
	check(err)

	org, err := SDKClient.GetUserProfileDefaultOrg(&profile[0])

	if err != nil {
		fmt.Println("Could not get default organization")
	} else {
		fmt.Println("Setting default org..")
		vpr.Set("org_id", org)
		vpr.WriteConfig()
	}
}

func setDefaultProviderKey(prov string) {
	var providerKey string

	if flagDebug {
		fmt.Printf("Debug - setDefaultProviderKey(%v)\n", prov)
	}

	switch prov {
	case "hci":
		providerKey = fmt.Sprintf("%s_keyset", prov)
	case "aws":
		providerKey = fmt.Sprintf("%s_keyset", prov)
	case "gce":
		providerKey = fmt.Sprintf("%s_keyset", prov)
	case "gke":
		providerKey = fmt.Sprintf("%s_keyset", prov)
	case "azr":
		providerKey = fmt.Sprintf("%s_keyset", prov)
	case "eks":
		providerKey = fmt.Sprintf("%s_keyset", prov)
	case "user_ssh":
		providerKey = "ssh_keyset"
	default:
		fmt.Printf("Error - '%s' is not a known provider\n", prov)
		os.Exit(1)
	}

	profile, err := SDKClient.GetUserProfile()
	check(err)

	ks, err := ext.GetUserProfileKeysetID(&profile[0], vpr.GetInt("org_id"), prov)

	if err != nil {
		fmt.Println("Could not get default keyset")
	} else {
		if flagDebug {
			fmt.Printf("Debug - setDefaultProviderKey(%v) - Setting provider key %s to %d\n", prov, providerKey, ks)
		}
		vpr.Set(providerKey, ks)
		vpr.WriteConfig()
	}

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

	vpr.Set("ssh_keyset_id", v[0].ID)

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

	vpr.Set("ssh_keyset_id", v[0].ID)

}

func bootstrapConfigFile() {
	if flagDebug {
		fmt.Printf("Debug - bootstrapConfigFile()\n")
	}
	configBootStrap = true
	newConfigFile()
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(listConfigCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(initConfigCmd)

	configSetCmd.AddCommand(configSetTokenCmd)
	configSetCmd.AddCommand(configSetURLCmd)
	configSetCmd.AddCommand(configSetClusterCmd)
	configSetCmd.AddCommand(configSetOrganizationCmd)
	configSetCmd.AddCommand(configSetWorkspaceCmd)

	configSetClusterCmd.Flags().IntVarP(&flagclusterID, "id", "i", vpr.GetInt("cluster_id"), "ID of target cluster")
	e := configSetClusterCmd.MarkFlagRequired("id")
	check(e)

	configSetOrganizationCmd.Flags().IntVarP(&flagOrganizationID, "id", "i", 0, "ID of organization")
	e = configSetOrganizationCmd.MarkFlagRequired("id")
	check(e)

	configSetWorkspaceCmd.Flags().IntVarP(&flagWorkspaceID, "id", "i", 0, "ID of workspace")
	e = configSetWorkspaceCmd.MarkFlagRequired("id")
	check(e)

	initConfigCmd.PersistentFlags().BoolVarP(&flagSetDefaults, "set-defaults", "", true, "Configure default values if possible")
}
