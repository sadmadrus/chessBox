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

import "fmt"

// Move перемещает фигуру из from в to.
// Валидность хода не проверяется, но по итогам ход переходит к
// другому игроку, количество полуходов и ходов обновляется,
// обновляются данные о пешке en passant и рокировках.
// Рокировка и проведение пешки реализованы отдельными методами.
func (b *Board) Move(from, to Square) error {
	if from < 0 {
		return fmt.Errorf("%w: %v", errSquareNotExist, from)
	}
	if to < 0 {
		return fmt.Errorf("%w: %v", errSquareNotExist, to)
	}
	pc := b.brd[from]
	if (pc%2 == 1 && b.blk) || (pc%2 == 0 && !b.blk) {
		return fmt.Errorf("out of turn")
	}
	if b.ep > 0 && to == b.ep && pc < 3 { // может оказаться en passant
		switch pc {
		case WhitePawn:
			if b.ep >= 5*8 && b.brd[to-8] == BlackPawn {
				b.brd[to-8] = 0
			}
		case BlackPawn:
			if b.ep < 3*8 && b.brd[to+8] == WhitePawn {
				b.brd[to+8] = 0
			}
		}
	}
	b.ep = 0
	if b.blk {
		b.fm++
	}
	b.blk = !b.blk
	if pc != BlackPawn && pc != WhitePawn && b.brd[to] == 0 {
		b.hm++
	} else {
		b.hm = 0
	}
	switch from {
	case 0:
		b.RemoveCastling(WhiteQueenside)
	case 7:
		b.RemoveCastling(WhiteKingside)
	case 56:
		b.RemoveCastling(BlackQueenside)
	case 63:
		b.RemoveCastling(BlackKingside)
	}
	switch pc {
	case BlackKing:
		b.RemoveCastling(BlackKingside)
		b.RemoveCastling(BlackQueenside)
	case WhiteKing:
		b.RemoveCastling(WhiteKingside)
		b.RemoveCastling(WhiteQueenside)
	case BlackPawn:
		if from >= 8*6 && from < 8*7 && to >= 8*4 && to < 8*5 {
			b.ep = from - 8
		}
	case WhitePawn:
		if from >= 8*1 && from < 8*2 && to >= 8*3 && to < 8*4 {
			b.ep = from + 8
		}
	}
	b.brd[from] = 0
	b.brd[to] = pc
	return nil
}

// Castle производит рокировку.
// Проверяется доступность (фигуры не двигались, пространство свободно),
// но не валидность (проверки на нахождение короля под шахом и/или
// движение через поле под шахом нет).
func (b *Board) Castle(c Castling) error {
	er := fmt.Errorf("the deserved castling is impossible")
	if c != BlackKingside && c != WhiteKingside && c != BlackQueenside && c != WhiteQueenside {
		return er
	}
	if !b.HaveCastling(c) {
		return er
	}
	row := 0
	k := WhiteKing
	if c == BlackKingside || c == BlackQueenside {
		row = 7
		k = BlackKing
	}

	// под 960 эту логику надо будет переписать
	kingFrom := 4
	rookFrom := 7
	rookTo := 5
	kingTo := 6
	mm := func(i int) int { return i + 1 }
	if c == BlackQueenside || c == WhiteQueenside {
		rookFrom = 0
		rookTo = 3
		kingTo = 2
		mm = func(i int) int { return i - 1 }
	}
	for i := rookTo; i != rookFrom; i = mm(i) {
		if b.brd[row*8+i] != 0 {
			return er
		}
	}

	if row == 0 {
		b.RemoveCastling(WhiteKingside)
		b.RemoveCastling(WhiteQueenside)
	} else {
		b.RemoveCastling(BlackKingside)
		b.RemoveCastling(BlackQueenside)
	}
	if err := b.Move(Sq(row*8+rookFrom), Sq(row*8+rookTo)); err != nil {
		return er
	}
	// Тут и ниже паника, а не возврат ошибки, потому что ладья уже переставлена,
	// и доска находится в неопределённом состоянии. Впрочем, если в коде нет
	// ошибок, эти паники в реальной жизни никогда не сработают; они тут для тестов.
	if err := b.Remove(Sq(row*8 + kingFrom)); err != nil {
		panic("this can not happen")
	}
	if err := b.Put(Sq(row*8+kingTo), k); err != nil {
		panic("this can not happen")
	}
	return nil
}

// Promote проводит пешку. Валидность хода не проверяется.
func (b *Board) Promote(from Square, to Square, p Piece) error {
	if from == -1 {
		return fmt.Errorf("%w: %v", errSquareNotExist, from)
	}
	if to == -1 {
		return fmt.Errorf("%w: %v", errSquareNotExist, to)
	}
	if p <= BlackPawn || p >= WhiteKing {
		return fmt.Errorf("can not promote to this piece")
	}
	if !(b.brd[from] == WhitePawn && from >= 6*8) && !(b.brd[from] == BlackPawn && from < 2*8) {
		return fmt.Errorf("promotion not possible")
	}
	if b.brd[from]%2 != p%2 {
		return fmt.Errorf("can not promote to wrong color")
	}
	if err := b.Move(from, to); err != nil {
		return err
	}
	b.brd[to] = p
	return nil
}
