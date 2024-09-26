package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "wechat",
	Short: "wechat bot server",
	Long:  `wechat bot server`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config-path", "", "config file (default is $HOME/wechat.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
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

		// Search config in home directory with name ".wechat" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("wechat")
	}
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
