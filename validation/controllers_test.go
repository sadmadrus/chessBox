package validation_test

import (
	"fmt"
	"github.com/sadmadrus/chessBox/internal/board"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestAdvanced(t *testing.T) {
	var (
		startBrd1WhiteUsFEN = "rnbq1bnr~ppP5~3p4~4pBBp~3PPPp1~QP2k1P1~P6P~R3K1NR+w+KQ+-+5+6"
		endBrd1_2WhiteUsFEN = "rnbq1bnr~ppP5~3p4~4pBBp~3PPPp1~QP2k1P1~P6P~2KR2NR+b+-+-+6+6"
		endBrd1_3WhiteUsFEN = "rBbq1bnr~pp6~3p4~4pBBp~3PPPp1~QP2k1P1~P6P~R3K1NR+b+KQ+-+5+6"
		endBrd1_4WhiteUsFen = "rnbq1bnr~ppP5~3Q4~4pBBp~3PPPp1~1P2k1P1~P6P~R3K1NR+b+KQ+-+5+6"

		startBrd1BlackUsFEN = "rnbq1bnr~ppP5~3p4~4pB1p~3PPPp1~QP2k1P1~P6P~R3K1NR+b+KQ+f3+5+6"
		endBrd1_2BlackUsFEN = "rnbq1bnr~ppP5~3p4~4pB1p~3PP3~QP2kpP1~P6P~R3K1NR+w+KQ+-+5+7"

		startBrd2BlackUsFEN = "rn2k2r~8~8~8~3q4~6n1~8~R3K2R+b+KQq+-+5+6"

		classicUsFen      = board.Classical().UsFEN()
		invalidboardUsFen = "dfahdfk 8jerq ~ + dfak"
	)

	srv := httptest.NewServer(http.HandlerFunc(validation.Advanced))
	defer srv.Close()

	tests := []struct {
		name         string
		brd          string
		from         string
		to           string
		newpiece     string
		resultStatus int
		resultBrd    string
	}{
		{"no piece", startBrd1WhiteUsFEN, "32", "33", "", 403, ""},
		{"pawn promotion successful", startBrd1WhiteUsFEN, "50", "57", "B", 200,
			fmt.Sprintf(`{"board":"%s"}`, endBrd1_3WhiteUsFEN)},
		{"white turn, black move", startBrd1WhiteUsFEN, "57", "40", "", 403, ""},
		{"Knight try diagonal move", startBrd1WhiteUsFEN, "6", "13", "", 403, ""},
		{"Q turn, P in the way", startBrd1WhiteUsFEN, "16", "18", "", 403, ""},
		{"B turn, p in the way", startBrd1WhiteUsFEN, "37", "23", "", 403, ""},
		{"P up, clash with p", startBrd1WhiteUsFEN, "22", "30", "", 403, ""},
		{"P up, clash with B", startBrd1WhiteUsFEN, "29", "37", "", 403, ""},
		{"R to P", startBrd1WhiteUsFEN, "0", "8", "", 403, ""},
		{"Q to p", startBrd1WhiteUsFEN, "16", "43", "", 200,
			fmt.Sprintf(`{"board":"%s"}`, endBrd1_4WhiteUsFen)},
		{"K O-O, N in the way", startBrd1WhiteUsFEN, "4", "6", "", 403, ""},
		{"K too close to k", startBrd1WhiteUsFEN, "4", "11", "", 403, ""},
		{"K O-O-O successful", startBrd1WhiteUsFEN, "4", "2", "", 200, fmt.Sprintf(`{"board":"%s"}`, endBrd1_2WhiteUsFEN)},

		{"p g4-f3 successful enPassant", startBrd1BlackUsFEN, "30", "21", "", 200,
			fmt.Sprintf(`{"board":"%s"}`, endBrd1_2BlackUsFEN)},
		{"p g4-h3 enPassant not allowed", startBrd1BlackUsFEN, "30", "23", "", 403, ""},
		{"k f3 under self-check", startBrd1BlackUsFEN, "20", "21", "", 403, ""},
		{"K O-O through checked cells", startBrd2BlackUsFEN, "4", "6", "", 403, ""},
		{"K O-O-O through checked cells", startBrd2BlackUsFEN, "4", "2", "", 403, ""},
		{"k O-O not allowed", startBrd2BlackUsFEN, "60", "62", "", 403, ""},
		{"k O-O-O through busy cells", startBrd2BlackUsFEN, "60", "58", "", 403, ""},

		{"e3-e2", classicUsFen, "e3", "e2", "", 403, ""},
		{"b2-c2", classicUsFen, "b2", "c2", "", 403, ""},
		{"K a1-c1", classicUsFen, "a1", "c1", "", 403, ""},
		{"k b7-b5", classicUsFen, "b7", "b5", "", 403, ""},
		{"Q e1-d3", classicUsFen, "e1", "d3", "", 403, ""},
		{"q a8-h7", classicUsFen, "a8", "h7", "", 403, ""},
		{"B c1-d3", classicUsFen, "c1", "d3", "", 403, ""},
		{"b c8-d8", classicUsFen, "c8", "d8", "", 403, ""},
		{"N b1-d1", classicUsFen, "b1", "d1", "", 403, ""},
		{"n b8-a7", classicUsFen, "b8", "a7", "", 403, ""},
		{"R 0-9", classicUsFen, "0", "9", "", 403, ""},
		{"r a8-b6", classicUsFen, "a8", "b6", "", 403, ""},
		{"e2-e4 with Q promotion", classicUsFen, "12", "28", "Q", 403, ""},
		{"e2-e4 with b promotion", classicUsFen, "12", "28", "b", 403, ""},
		{"no promotion, newpiece indicated", startBrd1WhiteUsFEN, "16", "24", "B", 403, ""},

		{"e2-e4 A", classicUsFen, "e2", "e4", "A", 400, ""},
		{"e2-e4 B1", classicUsFen, "e2", "e4", "B1", 400, ""},
		{"a8-a9", classicUsFen, "a8", "a9", "", 400, ""},
		{"b0-b1", classicUsFen, "b0", "b1", "", 400, ""},
		{"e2-e2", classicUsFen, "e2", "e2", "", 400, ""},
		{"62-67", classicUsFen, "62", "67", "", 400, ""},
		{"e2-", classicUsFen, "e2", "", "", 400, ""},
		{"invalid board", invalidboardUsFen, "e2", "e4", "", 400, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(srv.URL)
			v := u.Query()
			v.Add("board", tc.brd)
			v.Add("from", tc.from)
			v.Add("to", tc.to)
			v.Add("newpiece", tc.newpiece)
			u.RawQuery = v.Encode()

			res, err := http.Get(u.String())

			if err != nil {
				t.Fatal(err)
			}
			if tc.resultStatus != res.StatusCode {
				t.Fatalf("want %v, got %s", tc.resultStatus, res.Status)
			}

			var resBody []byte
			resBody, err = io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("error reading response body: %v", err)
			}
			if strings.TrimSpace(string(resBody)) != tc.resultBrd {
				t.Fatalf("JSON response: want %v, got %v", tc.resultBrd, strings.TrimSpace(string(resBody)))
			}
		})
	}
}

func FuzzAdvanced(f *testing.F) {
	srv := httptest.NewServer(http.HandlerFunc(validation.Advanced))
	defer srv.Close()
	tests := []string{
		"board=rnbq1bnr~ppP5~3p4~4pBBp~3PPPp1~QP2k1P1~P6P~R3K1NR+w+KQ+-+5+6&from=b3&to=b4&newpiece=",
		"to=b8&from=c7&newpiece=B&board=rnbq1bnr~ppP5~3p4~4pBBp~3PPPp1~QP2k1P1~P6P~R3K1NR+w+KQ+-+5+6",
		"from=0&to=7&board=rn2k2r~8~8~8~3q4~6n1~8~R3K2R+b+KQq+-+5+6",
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
