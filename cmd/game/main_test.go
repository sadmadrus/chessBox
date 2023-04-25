package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRoot404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(rootHandler))
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
	http.HandleFunc("/", rootHandler)
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
