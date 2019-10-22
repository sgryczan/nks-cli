package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	vpr "github.com/spf13/viper"
	ext "gitlab.com/sgryczan/nks-cli/nks/extensions"
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
	cobra.OnInitialize(initConfigSource, initClient, initCurrentConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nks.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&flagGenerateCompletionBash, "generatecompletion", "b", false, "Generate bash completion scripts")
	rootCmd.PersistentFlags().MarkHidden("generatecompletion")

	rootCmd.PersistentFlags().BoolVarP(&flagGenerateCompletionZsh, "generatecompletionzsh", "z", false, "Generate zsh completion scripts")
	rootCmd.PersistentFlags().MarkHidden("generatecompletionzsh")

	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "", false, "Debug logging")
	rootCmd.PersistentFlags().MarkHidden("debug")
}

func initClient() {
	if flagDebug {
		fmt.Printf("Debug - initClient()\n")
	}
	SDKClient = ext.NewClient(vpr.GetString("api_token"), vpr.GetString("api_url"))
}

// initConfig reads in config file and ENV variables if set.
func initConfigSource() {
	if flagDebug {
		fmt.Println("Debug - initConfigSource()")
	}

	// If a completion flag is present skip initialization
	if flagGenerateCompletionBash {
		rootCmd.GenBashCompletion(os.Stdout)
		os.Exit(0)
	}
	if flagGenerateCompletionZsh {
		rootCmd.GenZshCompletion(os.Stdout)
		os.Exit(0)
	}

	// Initialize from config file
	if cfgFile != "" {
		// Use config file from the flag.
		vpr.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nks" (without extension).
		vpr.AddConfigPath(home)
		vpr.SetConfigName(".nks")
	}

	vpr.ReadInConfig()

	// Initialize from environment
	vpr.SetEnvPrefix("nks") // NKS_<whatever>
	for _, key := range configFields {
		vpr.BindEnv(key)
	}
	vpr.AutomaticEnv() // read in environment variables that match

	if flagDebug {
		fmt.Printf("DEBUG - vpr.settings from environment: %+v\n", vpr.AllSettings())
		vpr.AllSettings()
	}

}

func initCurrentConfig() {
	if flagDebug {
		fmt.Println("Debug - initCurrentConfig()")
	}

	// If a config file is found, read it in.
	if err := vpr.ReadInConfig(); err != nil {
		if flagDebug {
			fmt.Println("initCurrentConfig() - Could not find config file!")
		}
		newConfigFile()
		//bootstrapConfigFile()
	}
}
