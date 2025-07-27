package main

import (
	"fmt"
	"os"

	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

func main() {
	client := realdebrid.NewClient(os.Getenv("REALDEBRID_API_KEY"))

	downloads, err := client.GetDownloads(&realdebrid.GetDownloadRequest{
		// Page:  1,
		// Limit: 5,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- List :: ", len(downloads))
	for _, download := range downloads {
		fmt.Println("ID:", download.ID)
		fmt.Println("Filename:", download.Filename)
		fmt.Println("File Size:", download.FileSize)
		fmt.Println("Download Link:", download.Download)
		fmt.Println("Generated:", download.Generated)
		fmt.Println("---")
	}

	// Delete a id from download list
	// Replace with your actual download ID
	if err := client.DeleteDownload("download_id"); err != nil {
		fmt.Println("Error deleting download:", err)
	} else {
		fmt.Println("Download deleted successfully.")
	}
}
