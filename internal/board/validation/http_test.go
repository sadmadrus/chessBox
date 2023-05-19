// Copyright 2023 The chessBox Crew
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board/validation"
)

func FuzzValidator(f *testing.F) {
	srv := httptest.NewServer(http.HandlerFunc(validation.Validator))
	defer srv.Close()
	tests := []string{
		"board=rnbqkbnr~pppppppp~8~8~8~8~PPPPPPPP~RNBQKBNR+w+KQkq+-+0+1",
		"board=1P6~8~4r3~3k4~8~8~3K1Q2~8+w+-+-+0+1",
		"this=that",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, query string) {
		u, _ := url.Parse(srv.URL)
		u.RawQuery = query
		res1, err := http.Get(u.String())
		if err != nil {
			return
		}
		res2, err := http.Head(u.String())
		if err != nil {
			return
		}
		if res1.Status != res2.Status {
			t.Fatalf("HEAD and GET responses differ: %v and %v", res2.Status, res1.Status)
		}
		if res1.StatusCode != http.StatusOK && res1.StatusCode != http.StatusForbidden && res1.StatusCode != http.StatusBadRequest {
			t.Fatalf("unexpected response: %s", res1.Status)
		}
	})
}

func TestValidator(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(validation.Validator))
	defer srv.Close()
	tests := []struct {
		name  string
		query string
		want  int
	}{
		{"good", "board=rnbqkbnr~pppppppp~8~8~8~8~PPPPPPPP~RNBQKBNR+w+KQkq+-+0+1", 200},
		{"bad", "board=1P6~8~4r3~3k4~8~8~3K1Q2~8+w+-+-+0+1", 403},
		{"ugly", "this=that", 400},
	}
	t.Run("post", func(t *testing.T) {
		u, _ := url.Parse(srv.URL)
		res, err := http.Post(u.String(), "text/html", nil)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("want 405, got %v", res.StatusCode)
		}
	})
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(srv.URL)
			u.RawQuery = tc.query
			res, err := http.Get(u.String())
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.want {
				t.Fatalf("want %v, got %v", tc.want, res.StatusCode)
			}
		})
	}
}
