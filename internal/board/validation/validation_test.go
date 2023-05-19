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
	"sort"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/internal/board/validation"
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
		{"castling with moved rook", "4k2r/p2P1p1p/rp1bPPp1/2p1b1P1/3p4/7P/PPP2Q2/R3K2R b KQkq - 0 1", false},
		{"both kings in check", "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP3kP1/P6P/R3K1NR w KQ - 6 7", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := board.FromFEN(tc.fen)
			if err != nil {
				t.Fatal(err)
			}
			got := validation.IsLegal(*b)
			if got != tc.want {
				t.Fatalf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestCheckedBy(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		sq   interface{}
		want []string
	}{
		{"many attacks", "8/5Q2/6nb/3Np3/4K2r/6P1/8/2B2R2 w - - 0 1", "f4",
			[]string{"c1", "f1", "g3", "e4", "h4", "d5", "e5", "g6", "h6", "f7"}},
		{"k none", "8/5P2/1QN1k3/8/1b1q2R1/2np3B/2r5/4K3 w - - 5 6", "e6", []string{}},
		{"K none", "8/5P2/1QN1k3/8/1b1q2R1/2np3B/2r5/4K3 w - - 5 6", "e1", []string{}},
		{"check by B", "8/5P2/1QN1k3/8/1b1q4/3p3B/2r5/4K3 w - - 5 6", "e6", []string{"h3"}},
		{"check by b", "8/5P2/1QN1k3/8/1b1q4/3p3B/2r5/4K3 w - - 5 6", "e1", []string{"b4"}},
		{"check by q", "8/5P2/1Q2k3/8/1b1q2R1/2np3B/8/2r1K3 w - - 5 6", "e6", []string{"b6"}},
		{"check by R", "8/5P2/1Q2k3/8/1b1q2R1/2np3B/8/2r1K3 w - - 5 6", "e1", []string{"c1"}},
		{"check by N", "8/2N5/4k3/8/3K4/8/2n5/8 w - - 5 6", "e6", []string{"c7"}},
		{"check by n", "8/2N5/4k3/8/3K4/8/2n5/8 w - - 5 6", "d4", []string{"c2"}},
		{"check by p", "2B5/1P6/k7/8/6p1/7K/8/8 w - - 5 6", "h3", []string{"g4"}},
		{"no check by P", "2B5/1P6/k7/8/6p1/7K/8/8 w - - 5 6", "a6", []string{}},
		{"king no threat", "8/8/8/8/8/4k3/8/K7 w - - 5 6", "a1", []string{}},
		{"king up", "8/8/8/8/8/4k3/8/K7 w - - 5 6", "e2", []string{"e3"}},
		{"king left", "8/8/8/8/8/4k3/8/K7 w - - 5 6", "f3", []string{"e3"}},
		{"king diag", "8/8/8/8/8/4k3/8/K7 w - - 5 6", "d4", []string{"e3"}},
		{"knight doesn't reach", "8/8/8/8/2b5/1k6/3K4/1r3N2 w - - 5 6", 11, []string{}},
		{"knight does check", "k7/8/1NK5/8/8/8/8/8 b - - 5 6", 56, []string{"b6"}},
		{"queen is far", "3q4/8/8/8/8/8/3K4/8 w - - 5 6", 11, []string{"d8"}},
		{"q & r are blocked vertically", "3q4/8/3p4/3r4/3b4/8/3K4/8 w - - 5 6", 11, []string{}},
		{"q & r are blocked horizontally", "8/8/8/8/8/8/3K1RQq/8 w - - 5 6", 11, []string{}},
		{"Q horizontally", "8/8/8/8/8/1q1Q2k1/8/7K b - - 5 6", 22, []string{"d3"}},
		{"R is close", "8/8/8/8/8/6kR/8/7K b - - 5 6", 22, []string{"h3"}},
		{"R is blocked", "6k1/8/8/6N1/8/8/8/6R1 b - - 5 6", 62, []string{}},
		{"q far up-right", "7q/8/8/8/8/8/8/K7 w - - 5 6", 0, []string{"h8"}},
		{"b down-left", "8/1N6/8/3K4/8/1b6/8/7Q w - - 5 6", 35, []string{"b3"}},
		{"p down-left & down-right", "8/8/3K4/2p1p3/8/8/8/8 w - - 5 6", "d6", []string{}},
		{"P down-left", "8/8/8/3r3q/8/5k2/4P3/8 b - - 5 6", 21, []string{"e2"}},
		{"B & Q hidden, Ps don't check", "8/1B6/2P5/5P2/4k3/5b2/2R5/1Q6 b - - 5 6", 28, []string{}},
		{"B down-right", "8/5P2/1QN1k3/8/1b1q4/3p3B/2r5/4K3 b - - 5 6", 44, []string{"h3"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := board.FromFEN(tc.fen)
			if err != nil {
				t.Fatal(err)
			}
			var s board.Square
			if sq, ok := tc.sq.(int); ok {
				s = board.Sq(sq)
			} else {
				s = board.Sq(tc.sq.(string))
			}
			got := validation.CheckedBy(s, *b)
			fail := func() { t.Fatalf("want %v\ngot%v", tc.want, got) }
			if len(got) != len(tc.want) {
				fail()
			}
			sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
			for i := range got {
				if got[i].String() != tc.want[i] {
					fail()
				}
			}
		})
	}
}
