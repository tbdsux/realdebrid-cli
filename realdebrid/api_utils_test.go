package realdebrid

import (
	"testing"
)

func TestGetTime(t *testing.T) {
	client := NewClient("")

	time, err := client.GetTime()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Log("Current time from RealDebrid API:", time)

	if time == "" {
		t.Fatal("expected non-empty time string")
	}
}

func TestGetTimeISO(t *testing.T) {
	client := NewClient("")

	timeISO, err := client.GetTimeISO()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Log("Current ISO time from RealDebrid API:", timeISO)

	if timeISO == "" {
		t.Fatal("expected non-empty ISO time string")
	}
}
