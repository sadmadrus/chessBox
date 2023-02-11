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

package board_test

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestPackage(t *testing.T) {
	b := board.Classical()
	err := b.Move(board.Sq("e2"), board.Sq("e4"))
	if err != nil {
		t.Fatal(err)
	}
	if b.Equals(board.Classical()) {
		t.Fatal("can not be equal")
	}
	if !b.IsEnPassant(board.Sq("e3")) {
		t.Fatal("want e3 to be en passant")
	}
	_ = b.Move(board.Sq("c7"), board.Sq("c5"))
	if b.IsEnPassant(board.Sq("e3")) {
		t.Fatal("en passant should have changed")
	}
	b1, err := board.FromFEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	if err != nil {
		t.Fatal(err)
	}
	if !b.Equals(b1) {
		t.Fatal("those moves should have led to that position")
	}

	e := b.GetEnPassant().String()
	if e != "c6" {
		t.Fatalf("want c6, got %s", e)
	}
	_ = b.Move(board.Sq(6), board.Sq("f3"))
	want := "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
	got := b.FEN()
	if got != want {
		t.Fatalf("want: %s\ngot: %s\n", want, got)
	}
	p, err := b.Get(board.Sq("b8"))
	if err != nil {
		t.Fatal(err)
	}
	if p != board.BlackKnight {
		t.Fatalf("wrong piece, have %s", p.String())
	}
	if err = b.Remove(board.Sq(64)); err == nil {
		t.Fatal("should be an error")
	}
	if err = b.Move(board.Sq("e8"), board.Sq("x8")); err == nil {
		t.Fatal("should be an error")
	}
	if err = b.Remove(board.Sq(61)); err != nil {
		t.Fatal(err)
	}
	if err = b.Move(board.Sq("e8"), board.Sq("f8")); err != nil {
		t.Fatal(err)
	}
	if b.HaveCastling(board.CastlingBlackKingside) || b.HaveCastling(board.CastlingBlackQueenside) {
		t.Fatal("black castlings should have been revoked")
	}
}

func TestFromFENErrors(t *testing.T) {
	tests := map[string]string{
		"not all fields":      "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w c6 0 2",
		"who is to move?":     "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR x KQkq c6 0 2",
		"strange en passant":  "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d9 0 2",
		"strange halfmove":    "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 () 2",
		"strange move number": "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 two",
		"missing rows":        "rnbqkbnr/pp1ppppp/2p5/4P3/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"bad piece":           "rnbqkbnr/pp1ppDpp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
	}
	for name, fen := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := board.FromFEN(fen)
			if err == nil {
				t.Fatal("no error")
			}
			if b != nil {
				t.Fatal("board not nil")
			}
		})
	}
}

func TestCastlingString(t *testing.T) {
	tests := map[string]struct {
		str string
		err bool
	}{
		"KQkq":  {"KQkq", false},
		"Qk":    {"Qk", false},
		"extra": {"KKQQkkqq", true},
		"empty": {"-", false},
		"nil":   {"", true},
		"error": {"42", true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var b board.Board
			err := b.SetCastlingString(tc.str)
			if err != nil && !tc.err {
				t.Fatal(err)
			}
			if tc.err && err == nil {
				t.Fatalf("want error, have none")
			}
			if err == nil {
				got := b.CastlingString()
				if got != tc.str {
					t.Fatalf("want %s, got %s", tc.str, got)
				}
			}
		})
	}
}

func TestUsfen(t *testing.T) {
	u := "rnbqkbnr~pppppppp~8~8~8~8~PPPPPPPP~RNBQKBNR+w+KQkq+-+0+1"
	b, err := board.FromUsFEN(u)
	if err != nil {
		t.Fatal(err)
	}
	if !b.Equals(board.Classical()) {
		t.Fatal("should match exactly")
	}
	if u != board.Classical().UsFEN() {
		t.Fatal("should match")
	}
}

func FuzzFEN(f *testing.F) {
	tests := []string{
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 two",
	}
	for _, tc := range tests {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, fen string) {
		b, err := board.FromFEN(fen)
		if b == nil && err == nil {
			t.Fatal("nil Board and no error")
		}
		if err != nil {
			if b != nil {
				t.Fatal("have error but Board not nil")
			}
			return
		}
		got := b.FEN()
		if got != fen {
			t.Errorf("input: %s\noutput: %s", fen, got)
		}
	})
}

func TestMove(t *testing.T) {
	tests := map[string]struct {
		pos  string
		from string
		to   string
		want string
	}{
		"simple": {
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			"g1", "f3",
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		},
		"en passant White": {
			"rnbqkbnr/ppp2ppp/4p3/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
			"e5", "d6",
			"rnbqkbnr/ppp2ppp/3Pp3/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},
		"en passant Black": {
			"B4n2/3ppp2/1Np1k2N/2P1P3/2pPp2p/6B1/5P2/Q3R1K1 b - d3 0 1",
			"c4", "d3",
			"B4n2/3ppp2/1Np1k2N/2P1P3/4p2p/3p2B1/5P2/Q3R1K1 w - - 0 2",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, _ := board.FromFEN(tc.pos)
			if err := b.Move(board.Sq(tc.from), board.Sq(tc.to)); err != nil {
				t.Fatal(err)
			}
			got := b.FEN()
			if got != tc.want {
				t.Fatalf("want %s\ngot %s", tc.want, got)
			}
		})
	}
}
