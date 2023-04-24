package game

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func TestGameGet(t *testing.T) {
	srv := serveNewGame(t)
	defer srv.Close()

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
	if state.FEN != startingFEN {
		t.Fatalf("want starting position, got %s", state.FEN)
	}
}

func TestGameHead(t *testing.T) {
	srv := serveNewGame(t)
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, res.StatusCode)
	}
}

func TestGameMakeMove(t *testing.T) {
	srv := serveNewGame(t)
	defer srv.Close()

	data := url.Values{
		"player": []string{"white"},
		"move":   []string{"e2e4"},
	}

	req, err := http.NewRequest(http.MethodPut, srv.URL, strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("content", "application/x-www-form-urlencoded")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("want OK, got %s", res.Status)
	}
}

func serveNewGame(t *testing.T) *httptest.Server {
	t.Helper()
	g, err := new("none", "none", "none")
	if err != nil {
		t.Fatal(err)
	}
	return httptest.NewServer(g.handler())
}
