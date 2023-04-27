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
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
			if res.FEN != tc.wantFEN {
				t.Fatalf("want %s, got %s", tc.wantFEN, res.FEN)
			}

		})
	}
}

func TestRecreate(t *testing.T) {
	moves := []string{"d2d4", "g8f6", "c2c4", "d7d6", "g1f3", "b8d7", "b1c3", "e7e5", "e2e4", "g7g6",
		"c1e3", "f8g7", "d4e5", "d6e5", "h2h3", "c7c6", "d1d2", "d8e7", "e1c1", "e8g8",
		"d2d6", "e7d6", "d1d6"} // https://www.chessgames.com/perl/chessgame?gid=1006866
	var mm []fullMove
	whiteToMove := true
	for _, m := range moves {
		hm := parseMove(m)
		if whiteToMove {
			mm = append(mm, fullMove{white: hm})
		} else {
			mm[len(mm)-1].black = hm
		}
		whiteToMove = !whiteToMove
	}

	g, err := retrace(mm)
	if err != nil {
		t.Fatal(err)
	}
	want := "r1b2rk1/pp1n1pbp/2pR1np1/4p3/2P1P3/2N1BN1P/PP3PP1/2K2B1R b - - 0 12"
	if g.board.FEN() != want {
		t.Fatalf("want %s, got %s", want, g.board.FEN())
	}
}

func parseMove(uci string) Move {
	m, _ := ParseUCI(uci)
	return m
}
