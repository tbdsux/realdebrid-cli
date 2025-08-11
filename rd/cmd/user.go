/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show user information",
	Long:  `View user information of the linked Real Debrid api key.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := internal.GetApiKey()
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		client := realdebrid.NewClient(apiKey)
		userInfo, err := client.GetUser()
		if err != nil {
			cmd.PrintErrf("Error fetching user information: %v\n", err)
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"#", "Account Information"})
		t.AppendRows([]table.Row{
			{"ID", userInfo.ID},
			{"Username", userInfo.Username},
			{"Email", userInfo.Email},
			{"Points", userInfo.Points},
			{"Locale", userInfo.Locale},
			{"Avatar", userInfo.Avatar},
			{"Type", userInfo.Type},
			{"Premium", userInfo.Premium},
			{"Expiration", userInfo.Expiration},
		})
		t.SetStyle(table.StyleLight)
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
