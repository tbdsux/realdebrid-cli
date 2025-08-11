package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage your configuration settings for the realdebrid-cli application.`,
	Run: func(cmd *cobra.Command, args []string) {
		t := table.NewWriter()

		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Setting", "Value"})

		t.AppendRow(table.Row{"API Key", viper.GetString("apiKey")})
		t.AppendSeparator()
		t.AppendRow(table.Row{"Config File", viper.ConfigFileUsed()})
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
