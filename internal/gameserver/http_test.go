package gameserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func FuzzRoot(f *testing.F) {
	tests := []string{
		"white=this&black=that&notify=noway",
		"move=e2e4&player=white",
		"this=that",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	st := NewMemoryStorage()

	f.Fuzz(func(t *testing.T, a string) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(a))
		rr := httptest.NewRecorder()
		HandleRoot(st).ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			return
		}
		loc := rr.Result().Header.Get("location")
		if loc == "" {
			t.Fatal("empty or no Location header")
		}
	})
}

func FuzzGameAPI(f *testing.F) {
	st := NewMemoryStorage()

	tests := []string{
		"move=e2e4&player=white",
		"player=black&takeback=true",
		"player=blue&move=e7e5",
		"player=white&forfeit=true",
		"this=that",
	}
	for _, tc := range tests {
		f.Add(tc)
	}
	g := newGame(st)

	f.Fuzz(func(t *testing.T, a string) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(a))
		rr := httptest.NewRecorder()
		handler(st, g).ServeHTTP(rr, req)
		// этот фаззер ловит паники
	})
}

func TestRoot404(t *testing.T) {
	st := NewMemoryStorage()
	srv := httptest.NewServer(http.HandlerFunc(HandleRoot(st)))
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

func TestCreate(t *testing.T) {
	st := NewMemoryStorage()
	http.HandleFunc("/", HandleRoot(st))
	data := url.Values{}

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
}

func TestGameGet(t *testing.T) {
	st := NewMemoryStorage()
	g := newGame(st)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(st, g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	if rr.Result().Header.Get("content-type") != urlencoded {
		t.Fatal("URLencoded reply expected")
	}

	b, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	data, err := url.ParseQuery(string(b))
	if err != nil {
		t.Fatal(err)
	}
	fen := data.Get("position")
	if fen != startingFEN {
		t.Fatalf("want starting position, got %s", fen)
	}
}

func TestGameHead(t *testing.T) {
	st := NewMemoryStorage()
	g := newGame(st)
	req := httptest.NewRequest(http.MethodHead, "/", nil)
	rr := httptest.NewRecorder()
	handler(st, g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}
}

func TestGameMakeMove(t *testing.T) {
	st := NewMemoryStorage()
	g := newGame(st)

	data := url.Values{
		"player": []string{"white"},
		"move":   []string{"e2e4"},
	}

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data.Encode()))
	req.Header.Set("content", urlencoded)

	rr := httptest.NewRecorder()
	handler(st, g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	b, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	res, err := url.ParseQuery(string(b))
	if err != nil {
		t.Fatal(err)
	}
	fen := res.Get("position")

	want := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	if fen != want {
		t.Fatalf("want %s, got %s", want, fen)
	}

	data.Set("player", "black")
	data.Set("move", "d8d1")

	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data.Encode()))
	req.Header.Set("content", urlencoded)

	rr = httptest.NewRecorder()
	g, _ = st.LoadGame(g.ID)
	handler(st, g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("want %v, got %v", http.StatusForbidden, rr.Result().Status)
	}

	b, err = io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	res, err = url.ParseQuery(string(b))
	if err != nil {
		t.Fatal(err)
	}
	fen = res.Get("position")

	if fen != want {
		t.Fatalf("want %s, got %s", want, fen)
	}

}

func TestGameForfeit(t *testing.T) {
	st := NewMemoryStorage()
	g := newGame(st)

	data := url.Values{
		"player":  []string{"black"},
		"forfeit": []string{"true"},
	}

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data.Encode()))
	req.Header.Set("content", urlencoded)

	rr := httptest.NewRecorder()
	handler(st, g).ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("want %v, got %v", http.StatusOK, rr.Result().Status)
	}

	want := "1-0"
	b, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	res, err := url.ParseQuery(string(b))
	if err != nil {
		t.Fatal(err)
	}
	result := res.Get("result")
	if result != want {
		t.Fatalf("want %s, got %s", want, result)
	}
}
