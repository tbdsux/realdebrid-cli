package main

import (
	"fmt"
	"os"

	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

func main() {
	client := realdebrid.NewClient(os.Getenv("REALDEBRID_API_KEY"))

	link, err := client.UnrestricLink(&realdebrid.UnrestrictProps{
		Link: "https://www.dailymotion.com/video/x8snguw", // First Google result I found :)
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Unrestricted Link Details:")
	fmt.Println("ID:", link.ID)
	fmt.Println("Filename:", link.Filename)
	fmt.Println("File Size:", link.FileSize)
	fmt.Println("Link:", link.Link)
	fmt.Println("Host:", link.Host)
	fmt.Println("Chunks:", link.Chunks)
	fmt.Println("Download:", link.Download)
	fmt.Println("Streamable:", link.Streamable)

	fmt.Print("\n\n ==== ")

	// check

	check, err := client.UnrestrictCheck(&realdebrid.UnrestrictCheckProps{
		Link: "https://www.dailymotion.com/video/x8snguw",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Unrestrict Check Details:")
	fmt.Println("ID:", check.Host)
	fmt.Println("Link:", check.Link)
	fmt.Println("Filename:", check.Filename)

	fmt.Print("\n\n ==== ")

	// folder
	folder, err := client.UnrestrictFolder(
		"https://www.dailymotion.com/video/x8snguw",
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Unrestricted Folder Details:", folder)

}
