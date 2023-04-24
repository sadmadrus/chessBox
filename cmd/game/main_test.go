package main

import (
	"io"
	"net/http"
	"net/http/httptest"
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
