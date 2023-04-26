package game

import (
	"errors"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestGame(t *testing.T) {
	tests := map[string]struct {
		haveFEN    string
		haveStatus state
		request    request
		wantFEN    string
		wantErr    error
	}{
		"move_after_won": {
			startingFEN, blackWon,
			request{player: white, kind: makeMove, move: parseMove("e2e4")},
			startingFEN, errGameOver,
		},
		"move_out_of_turn": {
			startingFEN, ongoing,
			request{player: black, kind: makeMove, move: parseMove("g8f6")},
			startingFEN, errWrongTurn,
		},
		"valid_move": {
			startingFEN, ongoing,
			request{player: white, kind: makeMove, move: parseMove("e2e4")},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", nil,
		},
		"illegal_move": {
			startingFEN, ongoing,
			request{player: white, kind: makeMove, move: parseMove("a1a4")},
			startingFEN, errInvalidMove,
		},
		"forfeit_after_draw": {
			startingFEN, drawn,
			request{player: white, kind: forfeit, move: nil},
			startingFEN, errGameOver,
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

func parseMove(uci string) halfMove {
	m, _ := parseUCI(uci)
	return m
}
