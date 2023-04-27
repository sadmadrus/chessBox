package game

import (
	"errors"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func TestGame(t *testing.T) {
	tests := map[string]struct {
		haveFEN    string
		haveStatus state
		request    Request
		wantFEN    string
		wantErr    error
	}{
		"move_after_won": {
			startingFEN, blackWon,
			Request{Player: White, Kind: MakeMove, Move: parseMove("e2e4")},
			startingFEN, ErrGameOver,
		},
		"move_out_of_turn": {
			startingFEN, ongoing,
			Request{Player: Black, Kind: MakeMove, Move: parseMove("g8f6")},
			startingFEN, ErrWrongTurn,
		},
		"valid_move": {
			startingFEN, ongoing,
			Request{Player: White, Kind: MakeMove, Move: parseMove("e2e4")},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", nil,
		},
		"illegal_move": {
			startingFEN, ongoing,
			Request{Player: White, Kind: MakeMove, Move: parseMove("a1a4")},
			startingFEN, ErrInvalidMove,
		},
		"forfeit_after_draw": {
			startingFEN, drawn,
			Request{Player: White, Kind: Forfeit, Move: nil},
			startingFEN, ErrGameOver,
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
			res, err := g.Do(tc.request)
			if err != nil {
				t.Fatal(err)
			}
			if res.State.FEN != tc.wantFEN {
				t.Fatalf("want %s, got %s", tc.wantFEN, res.State.FEN)
			}
			if !errors.Is(res.Error, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, res.Error)
			}
		})
	}
}

func parseMove(uci string) Move {
	m, _ := ParseUCI(uci)
	return m
}
