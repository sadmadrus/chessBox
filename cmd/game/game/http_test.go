package game_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/sadmadrus/chessBox/cmd/game/game"
)

func TestCreatorFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(game.Creator))
	defer srv.Close()

	data := url.Values{
		"white": []string{"localhost:5566"},
		"black": []string{"127.0.0.1:8899"},
	}
	res, err := http.Post(srv.URL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %v", res.StatusCode)
	}
}

func TestCreatorSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(game.Creator))
	defer srv.Close()

	data := url.Values{
		"white":  []string{"localhost:5566"},
		"black":  []string{"127.0.0.1:8899"},
		"notify": []string{"127.0.0.1:8934"},
	}
	res, err := http.Post(srv.URL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %v", res.StatusCode)
	}
	loc := res.Header.Get("location")
	u, _ := url.Parse(srv.URL)
	u.Path = loc
	res, err = http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == 404 {
		t.Fatal("got 404, but the game is created")
	}
}
