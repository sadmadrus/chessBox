package validation_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sadmadrus/chessBox/validation"
)

func TestSimple(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(validation.Simple))
	defer srv.Close()

	tests := []struct {
		name  string
		piece string
		from  string
		to    string
		isOk  bool
	}{
		{"e2-e4", "P", "e2", "e4", true},
		{"a7-a6", "p", "a7", "a6", true},
		{"K e1-g1 O_O", "K", "e1", "g1", true},
		{"k e8-c8 O_O_O", "k", "e8", "c8", true},
		{"Q a1-h8", "Q", "a1", "h8", true},
		{"q b7-b2", "q", "b7", "b2", true},
		{"B c1-a3", "B", "c1", "a3", true},
		{"b f7-g8", "b", "f7", "g8", true},
		{"N b1-a3", "N", "b1", "a3", true},
		{"n e4-c3", "n", "e4", "c3", true},
		{"R 0-6", "R", "0", "6", true},
		{"r h8-b8", "r", "h8", "b8", true},

		{"e3-e2", "P", "e3", "e2", false},
		{"b2-c2", "p", "b2", "c2", false},
		{"K a1-c1", "K", "a1", "c1", false},
		{"k b7-b5", "k", "b7", "b5", false},
		{"Q e1-d3", "Q", "e1", "d3", false},
		{"q a8-h7", "q", "a8", "h7", false},
		{"B c1-d3", "B", "c1", "d3", false},
		{"b c8-d8", "b", "c8", "d8", false},
		{"N b1-d1", "N", "b1", "d1", false},
		{"n b8-a7", "n", "b8", "a7", false},
		{"R 0-9", "R", "0", "9", false},
		{"r a8-b6", "r", "a8", "b6", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(srv.URL)
			v := u.Query()
			v.Add("piece", tc.piece)
			v.Add("from", tc.from)
			v.Add("to", tc.to)
			u.RawQuery = v.Encode()

			res, err := http.Get(u.String())
			if err != nil {
				t.Fatal(err)
			}
			if tc.isOk && res.StatusCode != http.StatusOK {
				t.Fatalf("want OK, got %s", res.Status)
			}
			if !tc.isOk && res.StatusCode != http.StatusForbidden {
				t.Fatalf("want Forbidden, got %s", res.Status)
			}
		})
	}

	testsForBadRequest := []struct {
		name  string
		piece string
		from  string
		to    string
	}{
		{"A e2-e4", "A", "e2", "e4"},
		{"P1 e2-e4", "P1", "e2", "e4"},
		{"K a8-a9", "K", "a8", "a9"},
		{"k b0-b1", "k", "b0", "b1"},
		{"B c1-c1", "B", "c1", "c1"},
		{"N 62-67", "N", "62", "67"},
		{"R e2-", "R", "e2", ""},
	}

	for _, tc := range testsForBadRequest {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(srv.URL)
			v := u.Query()
			v.Add("piece", tc.piece)
			v.Add("from", tc.from)
			v.Add("to", tc.to)
			u.RawQuery = v.Encode()

			res, err := http.Get(u.String())
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != http.StatusBadRequest {
				t.Fatalf("want BadRequest, got %s", res.Status)
			}
		})
	}

}
