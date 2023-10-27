package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "hue-cli",
	Short: "Control Philips Hue via CLI",
	Long: `hue-cli is a command-line interface for controlling Philips Hue lights and bridges.

  To use hue-cli, you can run various subcommands to discover, register, list, and control your Hue devices.

  Examples:
  - Discover Hue Bridge on the local network: hue-cli discover
  - Register with a Hue Bridge: hue-cli register
  - List devices connected to the Hue Bridge: hue-cli list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Usage()
		if err != nil {
			fmt.Println("Failed to show usage:", err)
			return
		}
	},
}

func Execute() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
	}

	return rootCmd.Execute()
}
