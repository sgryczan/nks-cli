package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/viper"
)

var cfgFile string
var apiToken string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nks",
	Short: "A command line utility for NKS",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initClient)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nks.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&flagGenerateCompletions, "generatecompletion", "b", false, "Generate bash completion scripts")
	rootCmd.PersistentFlags().MarkHidden("generatecompletion")

	rootCmd.PersistentFlags().BoolVarP(&FlagDebug, "debug", "", false, "Debug logging")
	rootCmd.PersistentFlags().MarkHidden("debug")
}

func initClient() {
	SDKClient = nks.NewClient(viper.GetString("api_token"), viper.GetString("api_url"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if flagGenerateCompletions {
		rootCmd.GenBashCompletion(os.Stdout);
		os.Exit(0)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nks" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nks")
	}

	viper.SetEnvPrefix("nks") // NKS_<whatever>
	for _, key := range configFields {
		viper.BindEnv(key)
	}
	viper.AutomaticEnv() // read in environment variables that match
	err := viper.Unmarshal(CurrentConfig)
	check(err)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		err = viper.Unmarshal(CurrentConfig)
		check(err)
		//fmt.Printf("Found existing config: %+v", CurrentConfig)
	} else {
		fmt.Println("Could not find config file!")
		bootstrapConfigFile()
	}

	if CurrentConfig.OrgID == 0 {
		configureDefaultOrganization()
	}

	if CurrentConfig.SSHKeySetId == 0 {
		fmt.Println("Default SSH keys not set. Configuring...")
		setDefaultSSHKey()
	}

	if CurrentConfig.Provider == "" {
		configSet("provider", "gce")
	}

	if CurrentConfig.ProviderKeySetID == 0 {
		setDefaultProviderKey(CurrentConfig.Provider)
	}
}
