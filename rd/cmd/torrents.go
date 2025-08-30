/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	showDownloads "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/show_downloads"
	showTorrents "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/show_torrents"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var torrentsPageReq int
var torrentsLimitReq int
var noTorrentDownload bool

var torrentsCmd = &cobra.Command{
	Use:   "torrents",
	Short: "List all torrents",
	Long: `List all torrents

Show the list of your torrents. Selected item will
automatically be downloaded, set '--no-download' otherwise.
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

		selected, err := showTorrents.ShowTorrentsList(torrents, torrentsPageReq)
		if err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}

		// Show table file info
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"#", "Torrent Details"})
		t.AppendRows([]table.Row{
			{"ID", selected.ID},
			{"Filename", selected.Filename},
			{"Size", internal.ByteCountSI(selected.Bytes)},
			{"Status", selected.Status},
			{"Date Added", selected.Added},
		})
		t.SetStyle(table.StyleLight)
		t.Render()

		if noTorrentDownload {
			// --no-download set
			return
		}

		if len(selected.Links) == 0 {
			// expect > 1
			fmt.Println(" :: Empty links, selected cannot be downloaded yet.")
			return
		}

		fmt.Print("\n\n")
		if len(selected.Links) > 1 {
			fmt.Printf(" ::> Queueing %d files\n\n", len(selected.Links))
		}

		for _, item := range selected.Links {
			// Unrestrict torrent file
			unrestrict, err := showTorrents.HandleUnrestrictFileLink(item, rdClient)
			if err != nil {
				cmd.PrintErrf("Error: %v\n", err)
				return
			}

			if unrestrict.Quitting {
				return
			}

			// Show download details
			j := table.NewWriter()
			j.SetOutputMirror(os.Stdout)
			j.AppendHeader(table.Row{"#", "Download Details"})
			j.AppendRows([]table.Row{
				{"ID", unrestrict.Result.ID},
				{"Filename", unrestrict.Result.Filename},
				{"Size", internal.ByteCountSI(unrestrict.Result.FileSize)},
				{"Download", unrestrict.Result.Download},
				{"Type", unrestrict.Result.MimeType},
			})
			j.SetStyle(table.StyleLight)
			j.Render()

			// Start downloading
			fmt.Println("\n  Downloading ::", unrestrict.Result.Filename)

			output, err := showDownloads.DoDownloadFile(realdebrid.Download{
				ID:       unrestrict.Result.ID,
				Filename: unrestrict.Result.Filename,
				Download: unrestrict.Result.Download,
			})
			if err != nil {
				cmd.PrintErrf("Error: %v", err)
				return
			}

			if output.Quitting {
				if output.Fail {
					// err message has been printed probably
					return
				}

				fmt.Print(showDownloads.ShowFailDLMessage("Stopped download."))
				return
			}

			fmt.Print(showDownloads.ShowSuccessDLMessage("Download success!"))
		}
	},
}

func init() {
	rootCmd.AddCommand(torrentsCmd)

	torrentsCmd.Flags().IntVarP(&torrentsPageReq, "page", "p", 1, "Page to show")
	torrentsCmd.Flags().IntVarP(&torrentsLimitReq, "limit", "l", 100, "Number of items to show")
	torrentsCmd.Flags().BoolVar(&noTorrentDownload, "no-download", false, "Disables auto-download on selected torrent item")
}
