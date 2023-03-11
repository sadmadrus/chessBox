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

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{"B00", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1", true},
		{"endgame", "8/8/4r3/3k4/8/8/3K1Q2/8 w - - 0 1", true},
		{"two kings", "8/2k5/4r3/3k4/8/8/3K1Q2/8 w - - 0 1", false},
		{"no king", "8/8/4r3/3k4/8/8/5Q2/8 w - - 0 1", false},
		{"extra pawn", "rnbqkbnr/pppppppp/8/8/4P3/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1", false},
		{"P in 8", "1P6/8/4r3/3k4/8/8/3K1Q2/8 w - - 0 1", false},
		{"P in 1", "8/8/4r3/3k4/8/8/3K1Q2/2P5 w - - 0 1", false},
		{"wrong king in check", "rnbqkbnr/pp3ppp/4p3/2pp4/Q1PP4/2N5/PP2PPPP/R1B1KBNR w KQkq - 0 1", false},
		{"king in check", "rnbqkbnr/pp3ppp/4p3/2pp4/Q1PP4/2N5/PP2PPPP/R1B1KBNR b KQkq - 0 1", true},
		{"triple check", "8/1k3Q2/7b/2N1p1q1/7r/2K3P1/8/1RB5 b - - 0 1", false},
		{"rooks double check", "6R1/4R1k1/8/3BK3/7q/8/8/8 b - - 0 1", true},
		{"bishops double check", "1k6/B4Q2/6nb/3NB3/7r/6P1/6K1/5R2 b - - 0 1", false},
		{"pawn-and-knight check", "4R1r1/6k1/5P2/3BK2N/7q/8/8/8 b - - 0 1", false},
		{"e.p.", "4R1r1/4N1k1/8/2pBK3/7q/5P2/8/8 w - c6 0 1", true},
		{"e.p. impossible", "4R1r1/4N1k1/8/2p1K3/3P3q/5P2/3B4/8 b - d3 0 1", false},
		{"too many queens", "6r1/p2Ppp1p/1p1QPPp1/2p3Pk/3p4/7P/PPP2Q1K/8 w - - 0 1", false},
		{"too many rooks", "5rr1/p2Ppp1p/1p2PPp1/2p3Pk/3p4/7P/PPP2Q1K/1r6 w - - 0 1", false},
		{"too many same-field bishops", "5r2/p2Ppp1p/1p1bPPp1/2p1b1Pk/3p4/7P/PPP2Q1K/1r6 w - - 0 1", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := board.FromFEN(tc.fen)
			if err != nil {
				t.Fatal(err)
			}
			got := b.IsValid()
			if got != tc.want {
				t.Fatalf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestThreatsTo(t *testing.T) {
	tb, _ := board.FromFEN("8/5Q2/6nb/3Np3/4K2r/6P1/8/2B2R2 w - - 0 1")
	s := board.Sq("f4")
	want := []string{"f7", "f1", "h4", "h6", "c1", "d5", "g6", "e5", "g3", "e4"}
	got := tb.ThreatsTo(s)
	fail := func() { t.Fatalf("want %v\ngot%v", want, got) }
	if len(got) != len(want) {
		fail()
	}
	for i := range got {
		if got[i].String() != want[i] {
			fail()
		}
	}
}
