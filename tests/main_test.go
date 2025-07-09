package tests

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/health")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
