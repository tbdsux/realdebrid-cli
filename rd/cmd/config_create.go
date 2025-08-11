/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	initconfig "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/init_config"
)

var apiKey string
var force bool

// initCmd represents the create command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new configuration file",
	Long: `Creates a new config file for the CLI application.

The config file will be created in the following path: $HOME/.realdebrid-cli.yaml
If the file already exists, it will not be overwritten.
`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := initconfig.AskConfigForSetup()

		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if output.Quitting {
			// do nothing
			return
		}

		if !output.Success {
			fmt.Print("Unknown error has happened, should not happen!")
			return
		}

		apiKey := output.TextInput.Value()
		if apiKey == "" {
			return
		}

		// Setup config

		viper.Set("apiKey", apiKey)

		var cfgFile string = path.Join(os.Getenv("HOME"), ".realdebrid-cli.yaml")

		fmt.Println("Creating config file at:", cfgFile)

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			if err := viper.WriteConfigAs(cfgFile); err != nil {
				cmd.PrintErrf("Error writing config file: %v\n", err)
				return
			}
		}

		if force {
			if err := viper.WriteConfigAs(cfgFile); err != nil {
				cmd.PrintErrf("Error writing config file: %v\n", err)
				return
			}
		} else {
			if err := viper.SafeWriteConfigAs(cfgFile); err != nil {
				cmd.PrintErrf("Error writing config file: %v\n", err)
				return
			}
		}

		fmt.Println("Config file created successfully.")
	},
}

func init() {
	configCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite the config file if it exists")
}
