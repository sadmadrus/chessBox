package game

import (
	"errors"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestErrors(t *testing.T) {
	tests := map[string]struct {
		haveFEN    string
		haveStatus state
		request    request
		wantFEN    string
		wantErr    error
	}{
		"move_after_won": {
			board.Classical().FEN(),
			blackWon,
			request{player: white, kind: makeMove, move: simpleMove{board.Sq("e2"), board.Sq("e4")}},
			board.Classical().FEN(),
			errGameOver,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := board.FromFEN(tc.haveFEN)
			if err != nil {
				t.Fatal(err)
			}
			have := &game{status: tc.haveStatus, board: *b}
			g, err := start("none", "none", "none", have)
			if err != nil {
				t.Fatal(err)
			}
			res, err := requestWithTimeout(tc.request, g)
			if err != nil {
				t.Fatal(err)
			}
			if res.state.FEN != tc.wantFEN {
				t.Fatalf("want %s, got %s", tc.wantFEN, res.state.FEN)
			}
			if !errors.Is(res.err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, res.err)
			}
		})
	}
}
