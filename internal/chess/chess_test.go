package chess_test

import (
	"errors"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/internal/chess"
)

const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func TestMove(t *testing.T) {
	tests := map[string]struct {
		haveStatus chess.State
		player     chess.Player
		move       chess.Move
		wantFEN    string
		wantErr    error
	}{
		"move_after_won": {
			chess.BlackWon,
			chess.White, parseMove("e2e4"),
			startingFEN, chess.ErrGameOver,
		},
		"move_out_of_turn": {
			chess.Ongoing,
			chess.Black, parseMove("g8f6"),
			startingFEN, chess.ErrWrongTurn,
		},
		"valid_move": {
			chess.Ongoing,
			chess.White, parseMove("e2e4"),
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", nil,
		},
		"illegal_move": {
			chess.Ongoing,
			chess.White, parseMove("a1a4"),
			startingFEN, chess.ErrInvalidMove,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g := &chess.Game{State: tc.haveStatus, StartingPosition: *board.Classical()}
			err := g.MakeMove(tc.move, tc.player)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
			pos := g.CurrentPosition()
			if pos.FEN() != tc.wantFEN {
				t.Fatalf("want %s, got %s", tc.wantFEN, pos.FEN())
			}

		})
	}
}

func TestCurrentPosition(t *testing.T) {
	moves := []string{"d2d4", "g8f6", "c2c4", "d7d6", "g1f3", "b8d7", "b1c3", "e7e5", "e2e4", "g7g6",
		"c1e3", "f8g7", "d4e5", "d6e5", "h2h3", "c7c6", "d1d2", "d8e7", "e1c1", "e8g8",
		"d2d6", "e7d6", "d1d6"} // https://www.chessgames.com/perl/chessgame?gid=1006866

	mm := make([]chess.Move, len(moves))
	for i, s := range moves {
		m := parseMove(s)
		mm[i] = m
	}

	g := chess.Game{StartingPosition: *board.Classical(), Moves: mm}
	want := "r1b2rk1/pp1n1pbp/2pR1np1/4p3/2P1P3/2N1BN1P/PP3PP1/2K2B1R b - - 0 12"
	pos := g.CurrentPosition()
	if pos.FEN() != want {
		t.Fatalf("want %s, got %s", want, pos.FEN())
	}
}

func TestForfeit(t *testing.T) {
	tests := map[string]struct {
		player chess.Player
		want   chess.State
	}{
		"white": {chess.White, chess.BlackWon},
		"black": {chess.Black, chess.WhiteWon},
		"other": {8, chess.Ongoing},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g := chess.Game{StartingPosition: *board.Classical()}
			_ = g.Forfeit(tc.player)
			if g.State != tc.want {
				t.Fatalf("want %v, got %v", tc.want, g.State)
			}
		})
	}
}

func parseMove(uci string) chess.Move {
	m, _ := chess.ParseUCI(uci)
	return m
}
