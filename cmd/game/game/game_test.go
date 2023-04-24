package game

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
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
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("want OK, got %s: %s", res.Status, string(b))
	}

	var state gameState
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}

	want := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	if state.FEN != want {
		t.Fatalf("want %s, got %s", want, state.FEN)
	}
}

func TestParseUCI(t *testing.T) {
	tests := []struct {
		move string
		from string
		to   string
	}{
		{"e2e4", "e2", "e4"},
	}

	for _, tc := range tests {
		t.Run(tc.move, func(t *testing.T) {
			m, err := parseUCI(tc.move)
			if err != nil {
				t.Fatal(err)
			}
			from := board.Sq(tc.from)
			to := board.Sq(tc.to)
			if from != m.fromSquare() || to != m.toSquare() {
				t.Fatalf("want from %v to %v, got from %v to %v", from, to, m.fromSquare(), m.toSquare())
			}
		})
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
