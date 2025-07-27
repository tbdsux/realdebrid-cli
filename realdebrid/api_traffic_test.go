package realdebrid

import (
	"os"
	"testing"
)

func TestGetTraffic(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	_, err := client.GetTraffic()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetTrafficDetails(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	// Use a valid start and end time for testing
	start := "2025-01-01"
	end := "2025-01-02"

	_, err := client.GetTrafficDetails(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
