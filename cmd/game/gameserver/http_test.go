package gameserver

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/sadmadrus/chessBox/internal/game"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const urlencoded = "application/x-www-form-urlencoded"

func TestRoot404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(RootHandler))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/does-not-exist")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %v", res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(b))
	want := "404 Game Not Found"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func TestCreateDelete(t *testing.T) {
	http.HandleFunc("/", RootHandler)
	data := url.Values{
		"white":  []string{"localhost:5566"},
		"black":  []string{"127.0.0.1:8899"},
		"notify": []string{"127.0.0.1:8934"},
	}

	// POST and check Location header
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data.Encode()))
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("want 201, got %v", rr.Result().Status)
	}
	loc := rr.Result().Header.Get("location")
	if loc == "" {
		t.Fatal("empty or no Location header")
	}

	// GET at Location
	req = httptest.NewRequest(http.MethodGet, loc, nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	// DELETE
	req = httptest.NewRequest(http.MethodDelete, loc, nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	// GET to verify the game is no longer there
	req = httptest.NewRequest(http.MethodGet, loc, nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("want %v, got %v", http.StatusNotFound, rr.Result().Status)
	}
}

func TestCreatorFail(t *testing.T) {
	data := url.Values{
		"white": []string{"localhost:5566"},
		"black": []string{"127.0.0.1:8899"},
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data.Encode()))
	rr := httptest.NewRecorder()
	http.HandlerFunc(creator).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %v", rr.Result().Status)
	}
}

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

	var state game.State
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
	req.Header.Set("content", urlencoded)

	rr := httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	var state game.State
	defer rr.Result().Body.Close()
	if err := json.NewDecoder(rr.Result().Body).Decode(&state); err != nil {
		t.Fatal(err)
	}

	want := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	if state.FEN != want {
		t.Fatalf("want %s, got %s", want, state.FEN)
	}

	data.Set("player", "black")
	data.Set("move", "d8d1")

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data.Encode()))
	req.Header.Set("content", urlencoded)

	rr = httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("want %v, got %v", http.StatusForbidden, rr.Result().Status)
	}

	defer rr.Result().Body.Close()
	if err := json.NewDecoder(rr.Result().Body).Decode(&state); err != nil {
		t.Fatal(err)
	}

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
	req.Header.Set("content", urlencoded)

	rr := httptest.NewRecorder()
	handler(g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	var state game.State
	defer rr.Result().Body.Close()
	if err := json.NewDecoder(rr.Result().Body).Decode(&state); err != nil {
		t.Fatal(err)
	}

	want := "1-0"
	if state.Status != want {
		t.Fatalf("want %s, got %s", want, state.Status)
	}
}

func serveNewGame(t *testing.T) game.ID {
	t.Helper()
	g, err := game.New("none", "none", "none")
	if err != nil {
		t.Fatal(err)
	}
	return g
}
