package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
  Use: "hue-cli",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("LOL")
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
