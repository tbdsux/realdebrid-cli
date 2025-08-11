package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

func getAvailableHosts(rd *realdebrid.RealDebridClient) {
	// === GET AVAILABLE HOSTS
	hosts, err := rd.GetTorrentsAvailableHosts()
	if err != nil {
		log.Println(err)
		return
	}

	for _, host := range hosts {
		fmt.Printf("Host: %s, Max File Size: %d\n", host.Host, host.MaxFileSize)
	}
}

func uploadTorrent(rd *realdebrid.RealDebridClient) string {
	// === UPLOAD A TORRENT
	cwd, _ := os.Getwd()
	torrentPath := path.Join(cwd, "examples", "torrent", "alice.torrent")

	res, err := rd.AddTorrent(torrentPath)
	if err != nil {
		log.Println(err)
		return ""
	}

	fmt.Println("ID: ", res.ID)
	fmt.Println("Upload URI: ", res.URI)

	return res.ID
}

func uploadMagnet(rd *realdebrid.RealDebridClient) string {
	// === Add Magnet
	// Same torrent as `alice.torrent` in example
	magnet := "magnet:?xt=urn:btih:OIX6MWZKUJWRJ423JLLCPUQCG3SIDWJE&dn=alice.txt&xl=163783"

	res, err := rd.AddMagnet(magnet)
	if err != nil {
		log.Println(err)
		return ""
	}

	fmt.Println("ID: ", res.ID)
	fmt.Println("Upload URI: ", res.URI)

	return res.ID
}

func selectTorrentFiles(rd *realdebrid.RealDebridClient, id string) {
	// === SELECT TORRENT FILES

	err := rd.SelectTorrentFiles(id, []string{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Torrent files selected successfully.")
}

func deleteTorrentId(rd *realdebrid.RealDebridClient, id string) {
	// === DELETE TORRENT
	err := rd.DeleteTorrent(id)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Torrent deleted successfully.")
}

func getTorrentInfo(rd *realdebrid.RealDebridClient, id string) {
	// === GET TORRENT INFO
	torrent, err := rd.GetTorrentsInfo(id)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Torrent ID: %s\n", torrent.ID)
	fmt.Printf("Name: %s\n", torrent.Filename)
	fmt.Printf("Size: %d bytes\n", torrent.Bytes)
	fmt.Printf("Status: %s\n", torrent.Status)
	fmt.Printf("Files:\n")
	for _, v := range torrent.Files {
		fmt.Printf(" - %d (%s)\n", v.ID, v.Path)
	}
}

// UNCOMMENT SOME OF THE FUNCTIONS BELOW TO TEST THEM
func main() {
	rd := realdebrid.NewClient(os.Getenv("REALDEBRID_API_KEY"))

	getTorrentInfo(rd, "W6Z5SYI5JG3ZA")

	getAvailableHosts(rd)

	uploadTorrent(rd)

	uploadMagnet(rd)

	selectTorrentFiles(rd, "W6Z5SYI5JG3ZA") // replace with your torrent ID

	deleteTorrentId(rd, "LWNVJWV7J5SKE") // replace with your torrent ID
}
