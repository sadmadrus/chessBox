package validation_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
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

			for _, method := range []string{"GET", "HEAD", "POST"} {
				client := &http.Client{}
				req, _ := http.NewRequest(method, u.String(), nil)
				res, err := client.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				switch method {
				case "GET", "HEAD":
					if tc.result != res.StatusCode {
						t.Fatalf("want %v, got %s", tc.result, res.Status)
					}

					if res.StatusCode == http.StatusOK {
						var resBody []byte
						resBody, err = io.ReadAll(res.Body)
						if err != nil {
							t.Fatalf("error reading response body: %v", err)
						}
						if string(resBody) != "" {
							t.Fatalf("JSON response: want no body response, got response %s", string(resBody))
						}
					}
				case "POST":
					if res.StatusCode != http.StatusMethodNotAllowed {
						t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
					}
				}
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

		for _, method := range []string{"GET", "HEAD", "POST"} {
			client := &http.Client{}
			req, _ := http.NewRequest(method, u.String(), nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			switch method {
			case "GET", "HEAD":
				if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusForbidden && res.StatusCode != http.StatusBadRequest {
					t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
				}

				if res.StatusCode == http.StatusOK {
					var resBody []byte
					resBody, err = io.ReadAll(res.Body)
					if err != nil {
						t.Fatalf("error reading response body: %v", err)
					}
					if string(resBody) != "" {
						t.Fatalf("JSON response: want no body response, got response %s", string(resBody))
					}
				}

			case "POST":
				if res.StatusCode != http.StatusMethodNotAllowed {
					t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
				}
			}
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

			for _, method := range []string{"GET", "HEAD", "POST"} {
				client := &http.Client{}
				req, _ := http.NewRequest(method, u.String(), nil)
				res, err := client.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				switch method {
				case "GET", "HEAD":
					if tc.resultStatus != res.StatusCode {
						t.Fatalf("want %v, got %s", tc.resultStatus, res.Status)
					}
				case "POST":
					if res.StatusCode != http.StatusMethodNotAllowed {
						t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
					}
				}

				if method == "GET" && res.StatusCode == http.StatusOK {
					var resBody []byte
					resBody, err = io.ReadAll(res.Body)
					if err != nil {
						t.Fatalf("error reading response body: %v", err)
					}
					if strings.TrimSpace(string(resBody)) != tc.resultBrd {
						t.Fatalf("JSON response: want %v, got %v", tc.resultBrd, strings.TrimSpace(string(resBody)))
					}
				}

				if method == "HEAD" && res.StatusCode == http.StatusOK {
					var resBody []byte
					resBody, err = io.ReadAll(res.Body)
					if err != nil {
						t.Fatalf("error reading response body: %v", err)
					}
					if string(resBody) != "" {
						t.Fatalf("JSON response: want no body response, got response %s", string(resBody))
					}
				}
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

		for _, method := range []string{"GET", "HEAD", "POST"} {
			client := &http.Client{}
			req, _ := http.NewRequest(method, u.String(), nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			switch method {
			case "GET", "HEAD":
				if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusForbidden && res.StatusCode != http.StatusBadRequest {
					t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
				}
			case "POST":
				if res.StatusCode != http.StatusMethodNotAllowed {
					t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
				}
			}

			if method == "GET" && res.StatusCode == http.StatusOK {
				var resBody []byte
				resBody, err = io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("error reading response body: %v", err)
				}
				if string(resBody) == "" {
					t.Fatal("JSON response: want body response, got no response")
				}
			}

			if method == "HEAD" && res.StatusCode == http.StatusOK {
				var resBody []byte
				resBody, err = io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("error reading response body: %v", err)
				}
				if string(resBody) != "" {
					t.Fatalf("JSON response: want no body response, got response %s", string(resBody))
				}
			}
		}
	})
}

func TestAvailableMoves(t *testing.T) {
	var (
		brdWhiteUsFEN = "3r4~p3PB2~2nr4~2k2Pp1~1b4Pq~1n1Q3P~1p1P1N2~RN2K2R+w+KQ+g6+5+6"
		brdBlackUsFEN = "3r4~p3PB2~2nr4~2k2Pp1~1b4Pq~1n1Q3P~1p1P1N2~RN2K2R+b+KQ+-+5+6"

		invalidboardUsFen = "dfahdfk 8jerq ~ + dfak"
	)

	srv := httptest.NewServer(http.HandlerFunc(validation.AvailableMoves))
	defer srv.Close()

	tests := []struct {
		name         string
		brd          string
		from         string
		resultStatus int
		resultMoves  string
	}{
		{"invalid board", invalidboardUsFen, "e2", 400, ""},
		{"not exist 65", brdWhiteUsFEN, "65", 400, ""},
		{"empty c1", brdWhiteUsFEN, "c1", 403, ""},

		{"R a1 w", brdWhiteUsFEN, "a1", 200, fmt.Sprintf(`{"moves":[8,16,24,32,40,48]}`)},
		{"R a1 b", brdBlackUsFEN, "a1", 403, ""},

		{"N b1 w", brdWhiteUsFEN, "1", 200, fmt.Sprintf(`{"moves":[16,18]}`)},
		{"N b1 b", brdBlackUsFEN, "1", 403, ""},

		{"K e1 w", brdWhiteUsFEN, "e1", 200, fmt.Sprintf(`{"moves":[3,5,6,12]}`)},
		{"K e1 b", brdBlackUsFEN, "e1", 403, ""},

		{"R h1 w", brdWhiteUsFEN, "7", 200, fmt.Sprintf(`{"moves":[5,6,15]}`)},
		{"R h1 b", brdBlackUsFEN, "7", 403, ""},

		{"p b2 w", brdWhiteUsFEN, "b2", 403, ""},
		{"p b2 b", brdBlackUsFEN, "b2", 200, fmt.Sprintf(`{"moves":[0]}`)},

		{"P d2 w", brdWhiteUsFEN, "11", 403, ""},
		{"P d2 b", brdBlackUsFEN, "11", 403, ""},

		{"N f2 w", brdWhiteUsFEN, "f2", 403, ""},
		{"N f2 b", brdBlackUsFEN, "f2", 403, ""},

		{"n b3 w", brdWhiteUsFEN, "17", 403, ""},
		{"n b3 b", brdBlackUsFEN, "17", 200, fmt.Sprintf(`{"moves":[0,2,11,27,32]}`)},

		{"Q d3 w", brdWhiteUsFEN, "d3", 200,
			fmt.Sprintf(`{"moves":[5,10,12,17,18,20,21,22,26,27,28,33,35,40,43]}`)},
		{"Q d3 b", brdBlackUsFEN, "d3", 403, ""},

		{"P h3 w", brdWhiteUsFEN, "23", 403, ""},
		{"P h3 b", brdBlackUsFEN, "23", 403, ""},

		{"b b4 w", brdWhiteUsFEN, "b4", 403, ""},
		{"b b4 b", brdBlackUsFEN, "b4", 200, fmt.Sprintf(`{"moves":[11,16,18,32]}`)},

		{"P g4 w", brdWhiteUsFEN, "30", 403, ""},
		{"P g4 b", brdBlackUsFEN, "30", 403, ""},

		{"q h4 w", brdWhiteUsFEN, "h4", 403, ""},
		{"q h4 b", brdBlackUsFEN, "h4", 200,
			fmt.Sprintf(`{"moves":[13,22,23,30,39,47,55,63]}`)},

		{"k c5 w", brdWhiteUsFEN, "34", 403, ""},
		{"k c5 b", brdBlackUsFEN, "34", 200, fmt.Sprintf(`{"moves":[41]}`)},

		{"P f5 w", brdWhiteUsFEN, "f5", 200, fmt.Sprintf(`{"moves":[45,46]}`)},
		{"P f5 b", brdBlackUsFEN, "f5", 403, ""},

		{"p g5 w", brdWhiteUsFEN, "38", 403, ""},
		{"p g5 b", brdBlackUsFEN, "38", 403, ""},

		{"n c6 w", brdWhiteUsFEN, "c6", 403, ""},
		{"n c6 b", brdBlackUsFEN, "c6", 200,
			fmt.Sprintf(`{"moves":[27,32,36,52,57]}`)},

		{"r d6 w", brdWhiteUsFEN, "43", 403, ""},
		{"r d6 b", brdBlackUsFEN, "43", 200,
			fmt.Sprintf(`{"moves":[19,27,35,44,45,46,47,51]}`)},

		{"p a7 w", brdWhiteUsFEN, "a7", 403, ""},
		{"p a7 b", brdBlackUsFEN, "a7", 200, fmt.Sprintf(`{"moves":[32,40]}`)},

		{"P e7 w", brdWhiteUsFEN, "52", 200, fmt.Sprintf(`{"moves":[59,60]}`)},
		{"P e7 b", brdBlackUsFEN, "52", 403, ""},

		{"B f7 w", brdWhiteUsFEN, "f7", 200,
			fmt.Sprintf(`{"moves":[17,26,35,39,44,46,60,62]}`)},
		{"B f7 b", brdBlackUsFEN, "f7", 403, ""},

		{"r d8 w", brdWhiteUsFEN, "59", 403, ""},
		{"r d8 b", brdBlackUsFEN, "59", 200,
			fmt.Sprintf(`{"moves":[51,56,57,58,60,61,62,63]}`)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(srv.URL)
			v := u.Query()
			v.Add("board", tc.brd)
			v.Add("from", tc.from)
			u.RawQuery = v.Encode()

			for _, method := range []string{"GET", "HEAD", "POST"} {
				client := &http.Client{}
				req, _ := http.NewRequest(method, u.String(), nil)
				res, err := client.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				switch method {
				case "GET", "HEAD":
					if tc.resultStatus != res.StatusCode {
						t.Fatalf("want %v, got %s", tc.resultStatus, res.Status)
					}
				case "POST":
					if res.StatusCode != http.StatusMethodNotAllowed {
						t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
					}
				}

				if method == "GET" && res.StatusCode == http.StatusOK {
					var resBody []byte
					resBody, err = io.ReadAll(res.Body)
					if err != nil {
						t.Fatalf("error reading response body: %v", err)
					}
					if strings.TrimSpace(string(resBody)) != tc.resultMoves {
						t.Fatalf("JSON response: want %v, got %v", tc.resultMoves, strings.TrimSpace(string(resBody)))
					}
				}

				if method == "HEAD" && res.StatusCode == http.StatusOK {
					var resBody []byte
					resBody, err = io.ReadAll(res.Body)
					if err != nil {
						t.Fatalf("error reading response body: %v", err)
					}
					if string(resBody) != "" {
						t.Fatalf("JSON response: want no body response, got response %s", string(resBody))
					}
				}
			}
		})
	}
}

func FuzzAvailableMoves(f *testing.F) {
	srv := httptest.NewServer(http.HandlerFunc(validation.AvailableMoves))
	defer srv.Close()
	tests := []string{
		"board=3r4~p3PB2~2nr4~2k2Pp1~1b4Pq~1n1Q3P~1p1P1N2~BN2K2R+w+KQ+-+5+6&from=b3",
		"from=8&board=rnbqkbnr~pppppppp~8~8~8~8~PPPPPPPP~RNBQKBNR+w+KQkq+-+1+1",
		"from=c8&board=rn2k2r~8~8~8~3q4~6n1~8~R3K2R+b+KQq+-+5+6",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, query string) {
		u, _ := url.Parse(srv.URL)
		u.RawQuery = query

		for _, method := range []string{"GET", "HEAD", "POST"} {
			client := &http.Client{}
			req, _ := http.NewRequest(method, u.String(), nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			switch method {
			case "GET", "HEAD":
				if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusForbidden && res.StatusCode != http.StatusBadRequest {
					t.Fatalf("unexpected reply for method %s: %s", method, res.Status)
				}
			case "POST":
				if res.StatusCode != http.StatusMethodNotAllowed && res.StatusCode != http.StatusBadRequest {
					t.Fatalf("unexpected reply for method %s: %s %v", method, res.Status, req)
				}
			}

			if method == "GET" && res.StatusCode == http.StatusOK {
				var resBody []byte
				resBody, err = io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("error reading response body: %v", err)
				}
				if string(resBody) == "" {
					t.Fatal("JSON response: want body response, got no response")
				}
			}

			if method == "HEAD" && res.StatusCode == http.StatusOK {
				var resBody []byte
				resBody, err = io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("error reading response body: %v", err)
				}
				if string(resBody) != "" {
					t.Fatalf("JSON response: want no body response, got response %s", string(resBody))
				}
			}
		}
	})
}
