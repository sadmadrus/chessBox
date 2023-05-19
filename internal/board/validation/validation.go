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

// Пакет validation содержит логику, анализирующую позицию на доске.
package validation

import "github.com/sadmadrus/chessBox/internal/board"

// IsLegal проверяет, могла ли позиция на доске легально возникнуть в ходе игры.
//
// Проверки не доскональные, ряд нелегальных позиций могут быть определены как
// легальные. А вот наоборот (чтоб легальная позиция была определена как
// нелегальная) случиться не должно.
func IsLegal(b board.Board) bool {
	bKing := board.Square(-1)
	wKing := board.Square(-1)
	var bPawns, wPawns []board.Square
	number := make(map[board.Piece]int)
	for i := 0; i < 64; i++ {
		s := board.Sq(i)
		c := getPiece(b, s)

		if c != 0 {
			number[c]++
		}
		switch c {
		case board.BlackKing:
			if bKing != -1 {
				return false
			}
			bKing = s
		case board.WhiteKing:
			if wKing != -1 {
				return false
			}
			wKing = s
		case board.BlackPawn:
			if in1(s) || in8(s) {
				return false
			}
			bPawns = append(bPawns, s)
		case board.WhitePawn:
			if in1(s) || in8(s) {
				return false
			}
			wPawns = append(wPawns, s)
		case board.BlackBishop, board.WhiteBishop:
			var pc board.Piece
			if s.IsBlack() {
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
	if !b.NextToMove() {
		thatKing, thisKing = thisKing, thatKing
	}
	if len(CheckedBy(thatKing, b)) > 0 {
		return false
	}
	if !checkCombinationLegal(b, CheckedBy(thisKing, b)) {
		return false
	}

	if !enPassantValid(b) {
		return false
	}

	if len(bPawns) > 8 || len(wPawns) > 8 {
		return false
	}

	if !numPiecesOk(number) {
		return false
	}

	if !castlingsValid(b) {
		return false
	}

	return true
}

// CheckedBy возвращает поля, на которых стоят фигуры, держащие данное поле «под
// боем». Если поле не пустое, в расчёт берутся только фигуры противоположного
// цвета.
func CheckedBy(s board.Square, b board.Board) []board.Square {
	if s < 0 || s > 63 {
		return nil
	}
	var isW, isB bool
	var out []board.Square
	if getPiece(b, s) != 0 {
		if getPiece(b, s)%2 == 0 {
			isB = true
		} else {
			isW = true
		}
	}

	// вертикали и горизонтали
	squares := make([]board.Square, 0, 4)
	for sq := s; sq < 64; sq += 8 {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
	}
	for sq := s; sq >= 0; sq -= 8 {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
	}
	for sq := s; ; sq++ {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inH(sq) {
			break
		}
	}
	for sq := s; ; sq-- {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inA(sq) {
			break
		}
	}
	for _, sq := range squares {
		if !isB && (getPiece(b, sq) == board.BlackQueen || getPiece(b, sq) == board.BlackRook) {
			out = append(out, sq)
		}
		if !isW && (getPiece(b, sq) == board.WhiteQueen || getPiece(b, sq) == board.WhiteRook) {
			out = append(out, sq)
		}
	}

	// диагонали
	squares = make([]board.Square, 0, 4)
	for sq := s; ; sq += 9 {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inH(sq) || in8(sq) {
			break
		}
	}
	for sq := s; ; sq += 7 {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inA(sq) || in8(sq) {
			break
		}
	}
	for sq := s; ; sq -= 7 {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inH(sq) || in1(sq) {
			break
		}
	}
	for sq := s; ; sq -= 9 {
		if getPiece(b, sq) != 0 && s != sq {
			squares = append(squares, sq)
			break
		}
		if inA(sq) || in1(sq) {
			break
		}
	}
	for _, sq := range squares {
		if !isB && (getPiece(b, sq) == board.BlackQueen || getPiece(b, sq) == board.BlackBishop) {
			out = append(out, sq)
		}
		if !isW && (getPiece(b, sq) == board.WhiteQueen || getPiece(b, sq) == board.WhiteBishop) {
			out = append(out, sq)
		}
	}

	// кони
	squares = make([]board.Square, 0, 8)
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
		if !isB && getPiece(b, sq) == board.BlackKnight {
			out = append(out, sq)
		}
		if !isW && getPiece(b, sq) == board.WhiteKnight {
			out = append(out, sq)
		}
	}

	// пешки
	if !in8(s) && !isB {
		squares = make([]board.Square, 0, 2)
		if !inA(s) {
			squares = append(squares, s+7)
		}
		if !inH(s) {
			squares = append(squares, s+9)
		}
		for _, sq := range squares {
			if getPiece(b, sq) == board.BlackPawn {
				out = append(out, sq)
			}
		}
	}
	if !in1(s) && !isW {
		squares = make([]board.Square, 0, 2)
		if !inA(s) {
			squares = append(squares, s-9)
		}
		if !inH(s) {
			squares = append(squares, s-7)
		}
		for _, sq := range squares {
			if getPiece(b, sq) == board.WhitePawn {
				out = append(out, sq)
			}
		}
	}

	// короли
	squares = make([]board.Square, 0, 8)
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
		if !isB && getPiece(b, sq) == board.BlackKing {
			out = append(out, sq)
		}
		if !isW && getPiece(b, sq) == board.WhiteKing {
			out = append(out, sq)
		}
	}

	return out
}

// inA указывает, находится ли поле на вертикали a.
func inA(s board.Square) bool {
	return s%8 == 0
}

// inH указывает, находится ли поле на вертикали h.
func inH(s board.Square) bool {
	return s%8 == 7
}

// in1 указывает, находится ли поле на горизонтали 1.
func in1(s board.Square) bool {
	return s >= 0 && s <= 7
}

// in8 указывает, находится ли поле на горизонтали 8.
func in8(s board.Square) bool {
	return s >= 56 && s <= 63
}

// checkCombinationLegal определяет, легально ли возник множественный шах.
func checkCombinationLegal(b board.Board, threats []board.Square) bool {
	if len(threats) > 2 {
		return false
	}
	if len(threats) == 2 {
		if getPiece(b, threats[0]) == getPiece(b, threats[1]) {
			if getPiece(b, threats[0]) < 7 {
				return false
			}
		}
		if getPiece(b, threats[0]) < 3 && getPiece(b, threats[1]) < 7 ||
			getPiece(b, threats[1]) < 3 && getPiece(b, threats[0]) < 7 {
			return false
		}
	}
	return true
}

// enPassantValid проверяет, похожа ли на правду информация о взятии на проходе.
func enPassantValid(b board.Board) bool {
	ep := b.GetEnPassant()
	if ep == -1 {
		return true
	}
	var p board.Piece
	var pSq, fromSq board.Square
	if !b.NextToMove() {
		if ep/8 != 2 {
			return false
		}
		p = board.WhitePawn
		pSq = ep + 8
		fromSq = ep - 8
	} else {
		if ep/8 != 5 {
			return false
		}
		p = board.BlackPawn
		pSq = ep - 8
		fromSq = ep + 8
	}

	return getPiece(b, pSq) == p && getPiece(b, ep) == 0 && getPiece(b, fromSq) == 0
}

// castlingsValid проверяет, на своих ли местах короли и ладьи для рокировок.
func castlingsValid(b board.Board) bool {
	if getPiece(b, board.Sq("e8")) != board.BlackKing && (b.HaveCastling(board.BlackKingside) || b.HaveCastling(board.BlackQueenside)) {
		return false
	}
	if getPiece(b, board.Sq("e1")) != board.WhiteKing && (b.HaveCastling(board.WhiteKingside) || b.HaveCastling(board.WhiteQueenside)) {
		return false
	}
	if getPiece(b, board.Sq("a1")) != board.WhiteRook && b.HaveCastling(board.WhiteQueenside) {
		return false
	}
	if getPiece(b, board.Sq("h1")) != board.WhiteRook && b.HaveCastling(board.WhiteKingside) {
		return false
	}
	if getPiece(b, board.Sq("a8")) != board.BlackRook && b.HaveCastling(board.BlackQueenside) {
		return false
	}
	if getPiece(b, board.Sq("h8")) != board.BlackRook && b.HaveCastling(board.BlackKingside) {
		return false
	}
	return true
}

// numPiecesOk проверяет, могло ли такое количество фигур на доске возникнуть не
// в Bughouse.
func numPiecesOk(num map[board.Piece]int) bool {
	var wExtras, bExtras int
	for p, n := range num {
		var e int
		switch p {
		case board.WhiteQueen, board.BlackQueen, 13, 14, 15, 16:
			e = n - 1
		case board.WhiteRook, board.BlackRook, board.WhiteKnight, board.BlackKnight:
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
	return wExtras+num[board.WhitePawn] <= 8 && bExtras+num[board.BlackPawn] <= 8
}

// getPiece возвращает фигуру в данной клетке, не проверяя валидность клетки.
func getPiece(b board.Board, s board.Square) board.Piece {
	p, _ := b.Get(s)
	return p
}
