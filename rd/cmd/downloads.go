package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	showDownloads "github.com/tbdsux/realdebrid-cli/rd/internal/handlers/show_downloads"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var downloadsPageReq int
var downloadsLimitReq int
var noDownload bool

var downloadsCmd = &cobra.Command{
	Use:   "downloads",
	Short: "List available downloadable files",
	Long: `List available downloadable files

Shows list of available files ready to be downloaded. Selected item will
automatically be downloaded, set '--no-download' otherwise.
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

		selected, err := showDownloads.ShowDownloadsList(downloads)
		if err != nil {
			cmd.PrintErrf("Error: %v", err)
			return
		}

		// Show table file info
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"#", "Download Details"})
		t.AppendRows([]table.Row{
			{"ID", selected.ID},
			{"Filename", selected.Filename},
			{"Size", internal.ByteCountSI(selected.FileSize)},
			{"Download", selected.Download},
			{"Type", selected.MimeType},
		})
		t.SetStyle(table.StyleLight)
		t.Render()

		if noDownload {
			// --no-download set
			return
		}

		fmt.Println("\n  Downloading ::", selected.Filename)

		output, err := showDownloads.DoDownloadFile(*selected)
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
	},
}

func init() {
	rootCmd.AddCommand(downloadsCmd)

	downloadsCmd.Flags().IntVarP(&downloadsPageReq, "page", "p", 1, "Page to show")
	downloadsCmd.Flags().IntVarP(&downloadsLimitReq, "limit", "l", 20, "Number of items to show")
	downloadsCmd.Flags().BoolVar(&noDownload, "no-download", false, "Disables auto-download on selected item")
}
