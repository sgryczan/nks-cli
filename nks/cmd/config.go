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

var configCmd = &cobra.Command{
	Use:   "config",
	Aliases: []string{"conf"},
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
	viper.WriteConfig()
	CurrentConfig.ClusterId = clusterId
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
	configSetCmd.AddCommand(configSetClusterCmd)

	configSetClusterCmd.Flags().IntVarP(&flagClusterId, "id", "i", CurrentConfig.ClusterId, "ID of target cluster")
	e := configSetClusterCmd.MarkFlagRequired("id")
	check(e)
}
