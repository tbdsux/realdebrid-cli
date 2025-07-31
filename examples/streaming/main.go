package main

import (
	"fmt"
	"os"

	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

func main() {
	client := realdebrid.NewClient(os.Getenv("REALDEBRID_API_KEY"))

	ID := "OW6O7PDIIPHEW"

	details, err := client.GetStreamingMediaInfos(ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Details ::>")
	fmt.Println("Filename:", details.Filename)
	fmt.Println("Hoster:", details.Hoster)
	fmt.Println("Link:", details.Link)
	fmt.Println("Type:", details.Type)
	fmt.Println("Season:", details.Season)
	fmt.Println("Episode:", details.Episode)
	fmt.Println("Year:", details.Year)
	fmt.Println("Duration:", details.Duration)
	fmt.Println("Bitrate:", details.Bitrate)

	transcode, err := client.GetStreamingTranscode(ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Transcode ::>", transcode)
}
