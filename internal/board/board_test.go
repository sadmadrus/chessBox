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
	b.Move(board.Sq("e2"), board.Sq("e4"))
	if b.Equals(board.Classical()) {
		t.Fatal("can not be equal")
	}
	if !b.IsEnPassant(board.Sq("e3")) {
		t.Fatal("want e3 to be en passant")
	}
	b.Move(board.Sq("c7"), board.Sq("c5"))
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
	b.Move(board.Sq(6), board.Sq("f3"))
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
			if _, err := board.FromFEN(fen); err == nil {
				t.Fatal("no error")
			}
		})
	}
}
