package shared

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

func TorrentDownload(torrent realdebrid.Torrent, rdClient *realdebrid.RealDebridClient, cmd *cobra.Command) {
	// Show table file info
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"#", "Torrent Details"})
	t.AppendRows([]table.Row{
		{"ID", torrent.ID},
		{"Filename", torrent.Filename},
		{"Size", internal.ByteCountSI(torrent.Bytes)},
		{"Status", torrent.Status},
		{"Date Added", torrent.Added},
	})
	t.SetStyle(table.StyleLight)
	t.Render()

	if torrent.Status != "downloaded" {
		// torrent cannot be downloaded yet
		fmt.Print(" ::> Torrent cannot be downloaded yet.\n\n")
		return
	}

	fmt.Print("\n\n")
	if len(torrent.Links) > 1 {
		fmt.Printf(" ::> Queueing %d files\n\n", len(torrent.Links))
	}

	for _, item := range torrent.Links {
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
}
