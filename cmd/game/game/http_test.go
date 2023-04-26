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
	g := serveNewGame(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	if rr.Result().Header.Get("content-type") != "application/json" {
		t.Fatal("JSON reply expected")
	}

	var state gameState
	defer rr.Result().Body.Close()
	if err := json.NewDecoder(rr.Result().Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	if state.FEN != startingFEN {
		t.Fatalf("want starting position, got %s", state.FEN)
	}
}

func TestGameHead(t *testing.T) {
	g := serveNewGame(t)
	req := httptest.NewRequest(http.MethodHead, "/", nil)
	rr := httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}
}

func TestGameMakeMove(t *testing.T) {
	g := serveNewGame(t)

	data := url.Values{
		"player": []string{"white"},
		"move":   []string{"e2e4"},
	}

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data.Encode()))
	req.Header.Set("content", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	var state gameState
	defer rr.Result().Body.Close()
	if err := json.NewDecoder(rr.Result().Body).Decode(&state); err != nil {
		t.Fatal(err)
	}

	want := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	if state.FEN != want {
		t.Fatalf("want %s, got %s", want, state.FEN)
	}
}

func TestGameForfeit(t *testing.T) {
	g := serveNewGame(t)

	data := url.Values{
		"player":  []string{"black"},
		"forfeit": []string{"true"},
	}

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data.Encode()))
	req.Header.Set("content", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	var state gameState
	defer rr.Result().Body.Close()
	if err := json.NewDecoder(rr.Result().Body).Decode(&state); err != nil {
		t.Fatal(err)
	}

	want := "1-0"
	if state.Status != want {
		t.Fatalf("want %s, got %s", want, state.Status)
	}
}

func serveNewGame(t *testing.T) id {
	t.Helper()
	g, err := start("none", "none", "none", nil)
	if err != nil {
		t.Fatal(err)
	}
	return g
}
