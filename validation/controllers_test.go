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
}
