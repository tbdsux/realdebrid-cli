package realdebrid

import (
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	user, err := client.GetUser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected non-zero user ID")
	}
	if user.Username == "" {
		t.Fatal("expected non-empty username")
	}
	if user.Email == "" {
		t.Fatal("expected non-empty email")
	}
}
