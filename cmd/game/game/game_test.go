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
			board.Classical().FEN(), blackWon,
			request{player: white, kind: makeMove, move: parseMove(t, "e2e4")},
			board.Classical().FEN(), errGameOver,
		},
		"move_out_of_turn": {
			board.Classical().FEN(), ongoing,
			request{player: black, kind: makeMove, move: parseMove(t, "g8f6")},
			board.Classical().FEN(), errWrongTurn,
		},
		"valid_move": {
			board.Classical().FEN(), ongoing,
			request{player: white, kind: makeMove, move: parseMove(t, "e2e4")},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", nil,
		},
		"illegal_move": {
			board.Classical().FEN(), ongoing,
			request{player: white, kind: makeMove, move: parseMove(t, "a1a4")},
			board.Classical().FEN(), errInvalidMove,
		},
		"forfeit_after_draw": {
			board.Classical().FEN(), drawn,
			request{player: white, kind: forfeit, move: nil},
			board.Classical().FEN(), errGameOver,
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

func parseMove(t *testing.T, uci string) halfMove {
	t.Helper()
	m, _ := parseUCI(uci)
	return m
}
