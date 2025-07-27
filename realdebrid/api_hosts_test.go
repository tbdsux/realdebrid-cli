package realdebrid

import (
	"os"
	"testing"
)

func TestGetHosts(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	_, err := client.GetHosts()
	if err != nil {
		t.Fatalf("Failed to get hosts: %v", err)
	}
}

func TestGetHostsStatus(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	_, err := client.GetHostsStatus()
	if err != nil {
		t.Fatalf("Failed to get hosts status: %v", err)
	}
}

func TestGetHostsRegex(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	_, err := client.GetHostsRegex()
	if err != nil {
		t.Fatalf("Failed to get hosts regex: %v", err)
	}
}

func TestGetHostsRegexFolder(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	_, err := client.GetHostsRegexFolder()
	if err != nil {
		t.Fatalf("Failed to get hosts regexFolder: %v", err)
	}
}

func TestGetHostsDomains(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	_, err := client.GetHostsDomains()
	if err != nil {
		t.Fatalf("Failed to get hosts domains: %v", err)
	}
}
