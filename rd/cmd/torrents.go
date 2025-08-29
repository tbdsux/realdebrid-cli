/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	showTorrents "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/show_torrents"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var torrentsPageReq int
var torrentsLimitReq int

var torrentsCmd = &cobra.Command{
	Use:   "torrents",
	Short: "List all torrents",
	Long: `List all torrents

Show the list of your torrents.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := internal.GetApiKey()
		if err != nil {
			cmd.PrintErrf("Error :%v", err)
			return
		}

		rdClient := realdebrid.NewClient(apiKey)

		torrents, err := rdClient.GetTorrents(&realdebrid.GetTorrentsRequest{
			Page:   torrentsPageReq,
			Limit:  torrentsLimitReq,
			Filter: "active",
		})
		if err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}

		if err := showTorrents.ShowTorrentsList(torrents, torrentsPageReq); err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(torrentsCmd)

	torrentsCmd.Flags().IntVarP(&torrentsPageReq, "page", "p", 1, "Page to show")
	torrentsCmd.Flags().IntVarP(&torrentsLimitReq, "limit", "l", 100, "Number of items to show")
}
