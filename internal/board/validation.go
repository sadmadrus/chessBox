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

// IsValid проверяет, могла ли позиция на доске легально возникнуть в ходе игры.
//
// Проверки не доскональные, ряд нелегальных позиций могут быть определены как
// легальные. А вот наоборот (чтоб легальная позиция была определена как
// нелегальная) случиться не должно.
func (b *Board) IsValid() bool {
	bKing := square(-1)
	wKing := square(-1)
	var bPawns, wPawns []square
	number := make(map[Piece]int)
	for s, c := range b.brd {
		if c != 0 {
			number[c]++
		}
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
			if in1(Sq(s)) || in8(Sq(s)) {
				return false
			}
			bPawns = append(bPawns, Sq(s))
		case WhitePawn:
			if in1(Sq(s)) || in8(Sq(s)) {
				return false
			}
			wPawns = append(wPawns, Sq(s))
		case BlackBishop, WhiteBishop:
			var pc Piece
			if Sq(s).isBlack() {
				pc = c + 8
			} else {
				pc = c + 10
			}
			number[pc]++
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
	if !b.checkCombinationLegal(b.ThreatsTo(thisKing)) {
		return false
	}

	if !b.enPassantValid() {
		return false
	}

	if len(bPawns) > 8 || len(wPawns) > 8 {
		return false
	}

	if !numPiecesOk(number) {
		return false
	}

	if !b.castlingsValid() {
		return false
	}

	return true
}

// ThreatsTo возвращает поля, на которых стоят фигуры, держащие данное поле «под
// боем». Если поле не пустое, в расчёт берутся только фигуры противоположного
// цвета.
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

	// вертикали и горизонтали
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

	// диагонали
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

	// кони
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

	// пешки
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

	// короли
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

// inA указывает, находится ли поле на вертикали a.
func inA(s square) bool {
	return s%8 == 0
}

// inH указывает, находится ли поле на вертикали h.
func inH(s square) bool {
	return s%8 == 7
}

// in1 указывает, находится ли поле на горизонтали 1.
func in1(s square) bool {
	return s >= 0 && s <= 7
}

// in8 указывает, находится ли поле на горизонтали 8.
func in8(s square) bool {
	return s >= 56 && s <= 63
}

// checkCombinationLegal определяет, легально ли возник множественный шах.
func (b *Board) checkCombinationLegal(threats []square) bool {
	if len(threats) > 2 {
		return false
	}
	if len(threats) == 2 {
		if b.brd[threats[0]] == b.brd[threats[1]] {
			if b.brd[threats[0]] < 7 {
				return false
			}
		}
		if b.brd[threats[0]] < 3 && b.brd[threats[1]] < 7 ||
			b.brd[threats[1]] < 3 && b.brd[threats[0]] < 7 {
			return false
		}
	}
	return true
}

// enPassantValid проверяет, похожа ли на правду информация о взятии на проходе.
func (b *Board) enPassantValid() bool {
	if b.ep == 0 {
		return true
	}
	var p Piece
	var pSq, fromSq square
	if b.blk {
		if b.ep/8 != 2 {
			return false
		}
		p = WhitePawn
		pSq = b.ep + 8
		fromSq = b.ep - 8
	} else {
		if b.ep/8 != 5 {
			return false
		}
		p = BlackPawn
		pSq = b.ep - 8
		fromSq = b.ep + 8
	}
	return b.brd[pSq] == p && b.brd[b.ep] == 0 && b.brd[fromSq] == 0
}

// castlingsValid проверяет, на своих ли местах короли и ладьи для рокировок.
func (b *Board) castlingsValid() bool {
	if b.brd[60] != BlackKing && (b.HaveCastling(BlackKingside) || b.HaveCastling(BlackQueenside)) {
		return false
	}
	if b.brd[4] != WhiteKing && (b.HaveCastling(WhiteKingside) || b.HaveCastling(WhiteQueenside)) {
		return false
	}
	if b.brd[0] != WhiteRook && b.HaveCastling(WhiteQueenside) {
		return false
	}
	if b.brd[7] != WhiteRook && b.HaveCastling(WhiteKingside) {
		return false
	}
	if b.brd[56] != BlackRook && b.HaveCastling(BlackQueenside) {
		return false
	}
	if b.brd[63] != BlackRook && b.HaveCastling(BlackKingside) {
		return false
	}
	return true
}

// numPiecesOk проверяет, могло ли такое количество фигур на доске возникнуть не
// в Bughouse.
func numPiecesOk(num map[Piece]int) bool {
	var wExtras, bExtras int
	for p, n := range num {
		var e int
		switch p {
		case WhiteQueen, BlackQueen, 13, 14, 15, 16:
			e = n - 1
		case WhiteRook, BlackRook, WhiteKnight, BlackKnight:
			e = n - 2
		}
		if e < 0 {
			continue
		}
		if p%2 == 1 {
			wExtras += e
		} else {
			bExtras += e
		}
	}
	return wExtras+num[WhitePawn] <= 8 && bExtras+num[BlackPawn] <= 8
}

// isBlack возвращает true, если поле чёрное.
func (s square) isBlack() bool {
	even := s%2 == 0
	if (s/8)%2 == 0 {
		return even
	}
	return !even
}
