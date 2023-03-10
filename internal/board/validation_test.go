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

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{"B00", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1", true},
		{"endgame", "8/8/4r3/3k4/8/8/3K1Q2/8 w - - 0 1", true},
		{"two kings", "8/2k5/4r3/3k4/8/8/3K1Q2/8 w - - 0 1", false},
		{"no king", "8/8/4r3/3k4/8/8/5Q2/8 w - - 0 1", false},
		{"extra pawn", "rnbqkbnr/pppppppp/8/8/4P3/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := board.FromFEN(tc.fen)
			got := b.IsValid()
			if got != tc.want {
				t.Fatalf("want %v, got %v", tc.want, got)
			}
		})
	}
}
