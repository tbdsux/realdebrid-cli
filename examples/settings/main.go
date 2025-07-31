package main

import (
	"fmt"
	"os"
	"path"

	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

func main() {
	client := realdebrid.NewClient(os.Getenv("REALDEBRID_API_KEY"))

	settings, err := client.GetSettings()
	if err != nil {
		panic(err)
	}

	fmt.Println("Settings ::>", settings)

	// Update user profile avatar
	cwd, _ := os.Getwd()
	profile := path.Join(cwd, "examples", "settings", "profile.png")

	if err := client.PutAvatarFile(profile); err != nil {
		panic(err)
	}

	fmt.Println("Avatar updated successfully.")
}
