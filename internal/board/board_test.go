package board_test

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestPackage(t *testing.T) {
	b, err := board.FromFEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	if err != nil {
		t.Fatal(err)
	}
	e := b.GetEnPassant().String()
	if e != "c6" {
		t.Fatalf("want c6, got %s", e)
	}
	b.Move(board.Sq(14), board.Sq("f3"))
}
