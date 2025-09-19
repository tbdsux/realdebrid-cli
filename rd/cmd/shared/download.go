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

	dirPath := ""
	if len(torrent.Links) > 1 {
		// Create directory for torrent files
		dirPath = torrent.Filename
		info, err := os.Stat(dirPath)
		if os.IsNotExist(err) {
			// Directory does not exist, create it
			if err := os.Mkdir(dirPath, 0755); err != nil {
				cmd.PrintErrf("Error creating directory: %v\n", err)
				return
			}
		} else {
			if !info.IsDir() {
				cmd.PrintErrf("Error: A file with the name '%s' already exists and is not a directory.\n", dirPath)
				return
			}
		}
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

		// Combine dir path and filename
		// cwd := os.Getwd()
		finalFilename := unrestrict.Result.Filename
		if dirPath != "" {
			finalFilename = fmt.Sprintf("%s/%s", dirPath, unrestrict.Result.Filename)
		}

		output, err := showDownloads.DoDownloadFile(realdebrid.Download{
			ID:       unrestrict.Result.ID,
			Filename: finalFilename,
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
