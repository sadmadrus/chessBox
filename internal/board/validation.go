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

package board

// IsValid attempts to check whether the position could legally appear during
// the game.
//
// The checks are not extremely thorough, so false positives are to be expected.
// There shouldn't be any false negatives, though.
func (b *Board) IsValid() bool {
	bKing := square(-1)
	wKing := square(-1)
	for s, c := range b.brd {
		switch c {
		case BlackKing:
			if bKing != -1 {
				return false
			}
			bKing = Sq(s)
		case WhiteKing:
			if wKing != -1 {
				return false
			}
			wKing = Sq(s)
		}
	}
	if bKing == -1 || wKing == -1 {
		return false
	}
	return true
}
