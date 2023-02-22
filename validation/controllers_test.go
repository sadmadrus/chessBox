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
		name   string
		piece  string
		from   string
		to     string
		result int
	}{
		{"e2-e4", "P", "e2", "e4", 200},
		{"a7-a6", "p", "a7", "a6", 200},
		{"K e1-g1 O_O", "K", "e1", "g1", 200},
		{"k e8-c8 O_O_O", "k", "e8", "c8", 200},
		{"Q a1-h8", "Q", "a1", "h8", 200},
		{"q b7-b2", "q", "b7", "b2", 200},
		{"B c1-a3", "B", "c1", "a3", 200},
		{"b f7-g8", "b", "f7", "g8", 200},
		{"N b1-a3", "N", "b1", "a3", 200},
		{"n e4-c3", "n", "e4", "c3", 200},
		{"R 0-6", "R", "0", "6", 200},
		{"r h8-b8", "r", "h8", "b8", 200},

		{"e3-e2", "P", "e3", "e2", 403},
		{"b2-c2", "p", "b2", "c2", 403},
		{"K a1-c1", "K", "a1", "c1", 403},
		{"k b7-b5", "k", "b7", "b5", 403},
		{"Q e1-d3", "Q", "e1", "d3", 403},
		{"q a8-h7", "q", "a8", "h7", 403},
		{"B c1-d3", "B", "c1", "d3", 403},
		{"b c8-d8", "b", "c8", "d8", 403},
		{"N b1-d1", "N", "b1", "d1", 403},
		{"n b8-a7", "n", "b8", "a7", 403},
		{"R 0-9", "R", "0", "9", 403},
		{"r a8-b6", "r", "a8", "b6", 403},

		{"A e2-e4", "A", "e2", "e4", 400},
		{"P1 e2-e4", "P1", "e2", "e4", 400},
		{"K a8-a9", "K", "a8", "a9", 400},
		{"k b0-b1", "k", "b0", "b1", 400},
		{"B c1-c1", "B", "c1", "c1", 400},
		{"N 62-67", "N", "62", "67", 400},
		{"R e2-", "R", "e2", "", 400},
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
			if tc.result != res.StatusCode {
				t.Fatalf("want %v, got %s", tc.result, res.Status)
			}
		})
	}
}

func FuzzSimple(f *testing.F) {
	srv := httptest.NewServer(http.HandlerFunc(validation.Simple))
	defer srv.Close()
	tests := []string{
		"piece=P&from=e2&to=e4",
		"from=4&to=2&piece=nil",
		"from=a7&to=b6&piece=Q",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, query string) {
		u, _ := url.Parse(srv.URL)
		u.RawQuery = query
		res, err := http.Get(u.String())
		if err != nil {
			return
		}
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusForbidden && res.StatusCode != http.StatusBadRequest {
			t.Fatalf("unexpected reply: %s", res.Status)
		}
	})
}
