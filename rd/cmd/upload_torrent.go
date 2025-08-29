package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	uploadtorrent "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/upload_torrent"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var torrentFile string
var autoSelectTorrent bool

// torrentCmd represents the torrent command
var torrentCmd = &cobra.Command{
	Use:   "upload-torrent",
	Short: "Upload a torrent file",
	Long: `Upload a torrent file

Select torrent file and upload to be downloaded later on.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if torrentFile == "" {
			return
		}

		if _, err := os.Stat(torrentFile); os.IsNotExist(err) {
			fmt.Println("[e] Torrent file does not exist:", torrentFile)
			return
		}

		if !strings.HasSuffix(torrentFile, ".torrent") {
			fmt.Println("[e] Invalid torrent file format. Please provide a .torrent file.")
			return
		}

		fmt.Println("[i] Torrent file provided:", torrentFile)

		apiKey, err := internal.GetApiKey()
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		rdClient := realdebrid.NewClient(apiKey)

		output, err := uploadtorrent.HandleUploadTorrent(
			torrentFile,
			rdClient,
		)

		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		if output.TaskDone {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"#", "Torrent Information"})
			t.AppendRow(table.Row{"File", output.TorrentFile})
			t.AppendRow(table.Row{"ID", output.Result.ID})
			t.AppendRow(table.Row{"URI", output.Result.URI})

			t.SetStyle(table.StyleLight)
			t.Render()
		}

		if autoSelectTorrent {
			fmt.Print("\n")

			// Do auto select and start torrent

			if err := uploadtorrent.AutoSelectFiles(output.Result.ID, rdClient); err != nil {
				cmd.PrintErrf("Error: %v\n", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(torrentCmd)

	torrentCmd.Flags().StringVarP(&torrentFile, "file", "f", "", "Path to the torrent file")
	torrentCmd.Flags().BoolVarP(&autoSelectTorrent, "auto-select", "a", false, "Automatically selects all the files to start the torrent")
	torrentCmd.MarkFlagRequired("file")
}
