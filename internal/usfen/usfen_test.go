package usfen_test

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/usfen"
)

const (
	classicFen   = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	classicUsfen = "rnbqkbnr~pppppppp~8~8~8~8~PPPPPPPP~RNBQKBNR+w+KQkq+-+0+1"
)

func TestFenToUsfen(t *testing.T) {
	got := usfen.FenToUsfen(classicFen)
	assert(t, classicUsfen, got)
}

func TestUsfenToFen(t *testing.T) {
	assert(t, classicFen, usfen.UsfenToFen(classicUsfen))
}

func assert[C comparable](t *testing.T, want, got C) {
	t.Helper()
	if want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
