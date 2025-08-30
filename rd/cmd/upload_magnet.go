package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/cmd/shared"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	addMagnet "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/add_magnet"
	uploadtorrent "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/upload_torrent"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var noAutoSelectMagnet bool
var doMagnetDL bool

var magnetCmd = &cobra.Command{
	Use:   "upload-magnet",
	Short: "Upload a magnet link",
	Long: `Upload a torrent magnet link

You will be asked to provide the magnet link on command usage.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := internal.GetApiKey()
		if err != nil {
			cmd.PrintErrf("Error :%v", err)
			return
		}

		inputMagnet, err := addMagnet.HandleAskMagnetLink()
		if err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}

		if inputMagnet.Quitting {
			return
		}

		magnetLink := inputMagnet.Textarea.Value()
		if magnetLink == "" {
			// empty
			return
		}

		rdClient := realdebrid.NewClient(apiKey)

		output, err := addMagnet.HandleUploadMagnetLink(magnetLink, rdClient)
		if err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}

		if output.TaskDone {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			t.AppendHeader(table.Row{"#", "Upload Details"})
			t.AppendRow(table.Row{"Magnet Link", fmt.Sprintf("%s...", output.MagnetLink[:60])})
			t.AppendRow(table.Row{"ID", output.Result.ID})
			t.AppendRow(table.Row{"URI", output.Result.URI})

			t.SetStyle(table.StyleLight)
			t.Render()
		}

		if noAutoSelectMagnet {
			return
		}

		fmt.Print("\n")

		// Do auto select and start torrent

		if err := uploadtorrent.AutoSelectFiles(output.Result.ID, rdClient); err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		if !doMagnetDL {
			return
		}

		fmt.Print("\n\n\n")

		// Get torrent
		torrent, err := rdClient.GetTorrentsInfo(output.Result.ID)
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		shared.TorrentDownload(*torrent.Torrent, rdClient, cmd)
	},
}

func init() {
	rootCmd.AddCommand(magnetCmd)

	magnetCmd.Flags().BoolVar(&noAutoSelectMagnet, "no-autoselect", false, "Disables auto selection of files once uploaded")
	magnetCmd.Flags().BoolVar(&doMagnetDL, "download", false, "Try to download once upload is complete")
}
