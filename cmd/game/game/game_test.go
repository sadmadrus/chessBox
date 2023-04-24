package game

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestGameGet(t *testing.T) {
	g, err := new("none", "none", "none")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(g.handler())
	res, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get("content-type") != "application/json" {
		t.Fatal("JSON reply expected")
	}

	var state gameState
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	if state.FEN != board.Classical().FEN() {
		t.Fatalf("want starting position, got %s", state.FEN)
	}
}
