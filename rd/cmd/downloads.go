/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	showDownloads "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/show_downloads"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var downloadsPageReq int
var downloadsLimitReq int

var downloadsCmd = &cobra.Command{
	Use:   "downloads",
	Short: "List available downloadable files",
	Long: `List available downloadable files

Shows list of available files ready to be downloaded.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := internal.GetApiKey()
		if err != nil {
			cmd.PrintErrf("Error :%v", err)
			return
		}

		rdClient := realdebrid.NewClient(apiKey)

		downloads, err := rdClient.GetDownloads(&realdebrid.GetDownloadRequest{
			Page:  downloadsPageReq,
			Limit: downloadsLimitReq,
		})
		if err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}

		if err := showDownloads.ShowDownloadsList(downloads); err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadsCmd)

	downloadsCmd.Flags().IntVarP(&downloadsPageReq, "page", "p", 1, "Page to show")
	downloadsCmd.Flags().IntVarP(&downloadsLimitReq, "limit", "l", 20, "Number of items to show")
}
