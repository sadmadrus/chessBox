package board_test

import (
	"fmt"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestSetFEN(t *testing.T) {
	b, err := board.FromFEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	fmt.Println(err)
	fmt.Println(*b)
}
