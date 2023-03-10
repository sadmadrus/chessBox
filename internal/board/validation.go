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
	var bPawns, wPawns []square
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
		case BlackPawn:
			if s < 8 || s > 55 {
				return false
			}
			bPawns = append(bPawns, Sq(s))
		case WhitePawn:
			if in1(Sq(s)) || in8(Sq(s)) {
				return false
			}
			wPawns = append(wPawns, Sq(s))
		}
	}
	if bKing == -1 || wKing == -1 {
		return false
	}

	thatKing := bKing
	thisKing := wKing
	if b.blk {
		thatKing, thisKing = thisKing, thatKing
	}
	if len(b.ThreatsTo(thatKing)) > 0 {
		return false
	}
	if len(b.ThreatsTo(thisKing)) > 2 {
		return false
	}

	if len(bPawns) > 8 || len(wPawns) > 8 {
		return false
	}
	return true
}

// ThreatsTo enumerates squares that the given square is threatened from.
// If the square is not empty, only enemy pieces/pawns are considered.
func (b *Board) ThreatsTo(s square) []square {
	if s < 0 || s > 63 {
		return nil
	}
	var isW, isB bool
	var out []square
	if b.brd[s] != 0 {
		if b.brd[s]%2 == 0 {
			isB = true
		} else {
			isW = true
		}
	}

	// vertical & horizontal
	squares := make([]square, 0, 4)
	for sq := s; sq < 64; sq += 8 {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
	}
	for sq := s; sq >= 0; sq -= 8 {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
	}
	for sq := s; ; sq++ {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inH(sq) {
			break
		}
	}
	for sq := s; ; sq-- {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inA(sq) {
			break
		}
	}
	for _, sq := range squares {
		if !isB && (b.brd[sq] == BlackQueen || b.brd[sq] == BlackRook) {
			out = append(out, sq)
		}
		if !isW && (b.brd[sq] == WhiteQueen || b.brd[sq] == WhiteRook) {
			out = append(out, sq)
		}
	}

	// diagonals
	squares = make([]square, 0, 4)
	for sq := s; ; sq += 9 {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inH(sq) || in8(sq) {
			break
		}
	}
	for sq := s; ; sq += 7 {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inA(sq) || in8(sq) {
			break
		}
	}
	for sq := s; ; sq -= 7 {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inH(sq) || in1(sq) {
			break
		}
	}
	for sq := s; ; sq -= 9 {
		if b.brd[sq] != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inA(sq) || in1(sq) {
			break
		}
	}
	for _, sq := range squares {
		if !isB && (b.brd[sq] == BlackQueen || b.brd[sq] == BlackBishop) {
			out = append(out, sq)
		}
		if !isW && (b.brd[sq] == WhiteQueen || b.brd[sq] == WhiteBishop) {
			out = append(out, sq)
		}
	}

	// knights
	squares = make([]square, 0, 8)
	if s%8 > 1 {
		if !in8(s) {
			squares = append(squares, s+6)
		}
		if !in1(s) {
			squares = append(squares, s-10)
		}
	}
	if !inA(s) {
		if s < 48 {
			squares = append(squares, s+15)
		}
		if s > 15 {
			squares = append(squares, s-17)
		}
	}
	if !inH(s) {
		if s < 48 {
			squares = append(squares, s+17)
		}
		if s > 15 {
			squares = append(squares, s-15)
		}
	}
	if s%8 < 6 {
		if !in8(s) {
			squares = append(squares, s+10)
		}
		if !in1(s) {
			squares = append(squares, s-6)
		}
	}
	for _, sq := range squares {
		if !isB && b.brd[sq] == BlackKnight {
			out = append(out, sq)
		}
		if !isW && b.brd[sq] == WhiteKnight {
			out = append(out, sq)
		}
	}

	// pawns
	if !in8(s) && !isB {
		squares = make([]square, 0, 2)
		if !inA(s) {
			squares = append(squares, s+7)
		}
		if !inH(s) {
			squares = append(squares, s+9)
		}
		for _, sq := range squares {
			if b.brd[sq] == BlackPawn {
				out = append(out, sq)
			}
		}
	}
	if !in1(s) && !isW {
		squares = make([]square, 0, 2)
		if !inA(s) {
			squares = append(squares, s-9)
		}
		if !inH(s) {
			squares = append(squares, s-7)
		}
		for _, sq := range squares {
			if b.brd[sq] == WhitePawn {
				out = append(out, sq)
			}
		}
	}

	// kings
	squares = make([]square, 0, 8)
	if !inA(s) {
		if !in1(s) {
			squares = append(squares, s-9)
		}
		squares = append(squares, s-1)
		if !in8(s) {
			squares = append(squares, s+7)
		}
	}
	if !inH(s) {
		if !in1(s) {
			squares = append(squares, s-7)
		}
		squares = append(squares, s+1)
		if !in8(s) {
			squares = append(squares, s+9)
		}
	}
	if !in1(s) {
		squares = append(squares, s-8)
	}
	if !in8(s) {
		squares = append(squares, s+8)
	}
	for _, sq := range squares {
		if !isB && b.brd[sq] == BlackKing {
			out = append(out, sq)
		}
		if !isW && b.brd[sq] == WhiteKing {
			out = append(out, sq)
		}
	}

	return out
}

// inA is true if square is on the "a" column.
func inA(s square) bool {
	return s%8 == 0
}

// inH is true if square is on the "h" column.
func inH(s square) bool {
	return s%8 == 7
}

// in1 is true if square is on the row 1.
func in1(s square) bool {
	return s >= 0 && s <= 7
}

// in8 is true if square is on the row 8.
func in8(s square) bool {
	return s >= 53 && s <= 63
}
