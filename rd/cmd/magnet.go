package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	addMagnet "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/add_magnet"
	uploadtorrent "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/upload_torrent"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var autoSelectMagnet bool

// magnetCmd represents the magnet command
var magnetCmd = &cobra.Command{
	Use:   "magnet",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

			t.AppendHeader(table.Row{"#", "Torrent Information"})
			t.AppendRow(table.Row{"Magnet Link", fmt.Sprintf("%s...", output.MagnetLink[:60])})
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
	rootCmd.AddCommand(magnetCmd)

	magnetCmd.Flags().BoolVarP(&autoSelectMagnet, "auto-select", "a", false, "Automatically selects all the files to start the torrent")
}
