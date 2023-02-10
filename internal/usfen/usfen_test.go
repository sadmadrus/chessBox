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
	got := usfen.FromFen(classicFen)
	assert(t, classicUsfen, got)
}

func TestUsfenToFen(t *testing.T) {
	assert(t, classicFen, usfen.ToFen(classicUsfen))
}

func assert[C comparable](t *testing.T, want, got C) {
	t.Helper()
	if want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
