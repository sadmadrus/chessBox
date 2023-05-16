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

// Пакет board реализует шахматную доску.
package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sadmadrus/chessBox/internal/usfen"
)

// Board представляет доску с позицией. Инициализация необязательна,
// пустой Board готов к употреблению, как пустая доска (без фигур,
// следующий ход первый).
type Board struct {
	brd [64]Piece // доска из 64 клеток
	blk bool      // false - ход белых, true - чёрных
	cas Castling  // битовая маска возможных рокировок
	ep  Square    // клетка, которая в прошлом ходу перепрыгнута пешкой
	hm  int       // полуходы без взятий и продвижения пешек
	fm  int       // номер хода -1 (для возможности использовать пустую доску)
}

// FromFEN возвращает доску из FEN-нотации. Валидность позиции не
// проверяется; на доске может оказаться 8 белых королей и чёрный ферзь,
// держащий под боем каждого короля, или белая пешка на 1-й горизонтали.
func FromFEN(fen string) (*Board, error) {
	ss := strings.Split(fen, " ")
	er := fmt.Errorf("%s is not a valid FEN", fen)
	if len(ss) != 6 {
		return nil, er
	}
	b := &Board{}
	switch ss[1] {
	case "w":
		b.blk = false
	case "b":
		b.blk = true
	default:
		return nil, er
	}
	if ss[3] != "-" {
		if ok := b.SetEnPassant(Sq(ss[3])); !ok {
			return nil, er
		}
	}
	hm, _ := strconv.Atoi(ss[4])
	if ss[4] != strconv.Itoa(hm) {
		return nil, er
	}
	fm, _ := strconv.Atoi(ss[5])
	if ss[5] != strconv.Itoa(fm) {
		return nil, er
	}
	b.hm = hm
	b.fm = fm

	if err := b.SetCastlingString(ss[2]); err != nil {
		return nil, er
	}

	ss = strings.Split(ss[0], "/")
	if len(ss) != 8 {
		return nil, er
	}
	for row := 0; row < 8; row++ {
		c := 0
		if len(ss[7-row]) == 0 {
			return nil, er
		}
		digit := false
		for _, r := range ss[7-row] {
			if c > 7 {
				return nil, er
			}
			switch r {
			case 'K':
				b.brd[row*8+c] = WhiteKing
			case 'k':
				b.brd[row*8+c] = BlackKing
			case 'Q':
				b.brd[row*8+c] = WhiteQueen
			case 'q':
				b.brd[row*8+c] = BlackQueen
			case 'R':
				b.brd[row*8+c] = WhiteRook
			case 'r':
				b.brd[row*8+c] = BlackRook
			case 'B':
				b.brd[row*8+c] = WhiteBishop
			case 'b':
				b.brd[row*8+c] = BlackBishop
			case 'N':
				b.brd[row*8+c] = WhiteKnight
			case 'n':
				b.brd[row*8+c] = BlackKnight
			case 'P':
				b.brd[row*8+c] = WhitePawn
			case 'p':
				b.brd[row*8+c] = BlackPawn
			default:
				if digit {
					return nil, er
				}
				v := int(r - '0')
				if v < 1 || v > 8 {
					return nil, er
				}
				c += v - 1
				digit = true
			}
			c++
			if r > '8' {
				digit = false
			}
		}
		if c != 8 {
			return nil, er
		}
	}
	return b, nil
}

// FromUsFEN возвращает доску из UsFEN-нотации. Валидность позиции не
// проверяется.
func FromUsFEN(s string) (*Board, error) {
	return FromFEN(usfen.ToFen(s))
}

// Classical возвращает доску, готовую к игре в классический
// вариант шахмат.
func Classical() *Board {
	b, _ := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return b
}

// Get возвращает фигуру, стоящую в заданной клетке.
func (b *Board) Get(s Square) (Piece, error) {
	if !s.IsValid() {
		return 0, fmt.Errorf("%w: %v", errSquareNotExist, s)
	}
	return b.brd[s], nil
}

// Put ставит фигуру на заданную клетку; возвращает ошибку,
// если клетка не пуста.
func (b *Board) Put(s Square, p Piece) error {
	if !s.IsValid() {
		return fmt.Errorf("%w: %v", errSquareNotExist, s)
	}
	if b.brd[s] != 0 {
		return fmt.Errorf("square not empty")
	}
	if p < 1 || p > BlackKing {
		return fmt.Errorf("unknown piece")
	}
	b.brd[s] = p
	return nil
}

// Remove убирает фигуру, стоящую в заданной клетке.
func (b *Board) Remove(s Square) error {
	if !s.IsValid() {
		return fmt.Errorf("%w: %v", errSquareNotExist, s)
	}
	if b.brd[s] == 0 {
		return fmt.Errorf("%v is empty", s)
	}
	b.brd[s] = 0
	return nil
}

// WhiteToMove устанавливает следующих ход белых.
func (b *Board) WhiteToMove() {
	b.blk = false
}

// BlackToMove устанавливает следующих ход чёрных.
func (b *Board) BlackToMove() {
	b.blk = true
}

// NextToMove возвращает true, если следующий ход белых,
// false, если следующий ход чёрных.
func (b *Board) NextToMove() bool {
	return !b.blk
}

// SetCastlingString устанавливает, какие рокировки доступны, по строке
// K, Q - королевская и ферзевая ладья для белых, k, q - для чёрных,
// "-" - рокировки недоступны.
func (b *Board) SetCastlingString(s string) error {
	err := fmt.Errorf("malformed castling string")
	have := b.cas
	b.cas = 0
	if s == "" {
		b.cas = have
		return err
	}
	if s == "-" {
		return nil
	}
	if strings.Contains(s, "K") {
		b.SetCastling(WhiteKingside)
	}
	if strings.Contains(s, "k") {
		b.SetCastling(BlackKingside)
	}
	if strings.Contains(s, "Q") {
		b.SetCastling(WhiteQueenside)
	}
	if strings.Contains(s, "q") {
		b.SetCastling(BlackQueenside)
	}
	if s != b.CastlingString() {
		b.cas = have
		return err
	}
	return nil
}

// CastlingString возвращает строку с перечислением возможных рокировок.
func (b *Board) CastlingString() string {
	if b.cas == 0 {
		return "-"
	}
	sb := strings.Builder{}
	if b.HaveCastling(WhiteKingside) {
		sb.WriteByte('K')
	}
	if b.HaveCastling(WhiteQueenside) {
		sb.WriteByte('Q')
	}
	if b.HaveCastling(BlackKingside) {
		sb.WriteByte('k')
	}
	if b.HaveCastling(BlackQueenside) {
		sb.WriteByte('q')
	}
	return sb.String()
}

// SetCastling устанавливает доступность рокировки.
func (b *Board) SetCastling(c Castling) {
	b.cas = b.cas | c
}

// HaveCastling проверяет доступность рокировки.
func (b *Board) HaveCastling(c Castling) bool {
	return b.cas&c != 0
}

// RemoveCastling убирает возможность рокировки.
func (b *Board) RemoveCastling(c Castling) {
	b.cas = b.cas &^ c
}

// SetEnPassant устанавливает клетку, перепрыгнутую пешкой в прошлом
// полуходу. Валидность не проверяется, только горизонталь // TODO?
func (b *Board) SetEnPassant(s Square) bool {
	if s/8 != 2 && s/8 != 5 {
		return false
	}
	b.ep = s
	return true
}

// GetEnPassant возвращает клетку, перепрыгнутую пешкой в прошлом ходу,
// -1 если такой не было.
func (b *Board) GetEnPassant() Square {
	if b.ep == 0 {
		return -1
	}
	return b.ep
}

// IsEnPassant возвращает true, если заданная клетка была перепрыгнута
// пешкой в прошлом ходу.
func (b *Board) IsEnPassant(s Square) bool {
	if s == 0 { // эта проверка нужна, потому что Board{} — валидная доска
		return false
	}
	return s == b.ep
}

// GetMoveNumber возвращает номер хода (полного).
func (b *Board) GetMoveNumber() int {
	return b.fm + 1
}

// SetMoveNumber устанавливает номер хода (полного) на доске.
func (b *Board) SetMoveNumber(n int) {
	b.fm = n - 1
}

// GetHalfMoves возвращает количество полуходов без взятия фигур и движения пешек.
func (b *Board) HalfMoves() int {
	return b.hm
}

// SetHalfMoves устанавливает количество полуходов без взятия фигур и движения пешек.
func (b *Board) SetHalfMoves(n int) {
	b.hm = n
}

// FEN возвращает FEN-нотацию доски.
func (b *Board) FEN() string {
	sb := strings.Builder{}
	for row := 7; row >= 0; row-- {
		for col := 0; col < 8; col++ {
			sb.WriteString(b.brd[8*row+col].String())
		}
		if row != 0 {
			sb.WriteByte('/')
		}
	}

	// strings.Replacer:
	// "--------" -> "8"
	// "-------" -> "7"
	// ...
	// "-" -> "1"
	rr := make([]string, 16)
	for i := 8; i > 0; i-- {
		rr[(8-i)*2] = strings.Repeat("-", i)
		rr[(8-i)*2+1] = strconv.Itoa(i)
	}
	r := strings.NewReplacer(rr...)
	s := r.Replace(sb.String())
	sb.Reset()
	sb.WriteString(s)
	if b.blk {
		sb.WriteString(" b ")
	} else {
		sb.WriteString(" w ")
	}
	sb.WriteString(b.CastlingString())
	sb.WriteRune(' ')
	sb.WriteString(b.GetEnPassant().String())
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(b.hm))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(b.fm))
	return sb.String()
}

// UsFEN возвращает UsFEN-нотацию доски.
func (b *Board) UsFEN() string {
	return usfen.FromFen(b.FEN())
}

// Equals указывает, одинаковая ли позиция на досках, с точностью до хода.
func (b *Board) Equals(that *Board) bool {
	return b.SamePosition(that) && b.hm == that.hm && b.fm == that.fm
}

// SamePosition возвращает true, если на обеих досках одинаковые позиции
// и очерёдность хода (номер хода может не совпадать - например, позиция
// возникла в той же партии повторно).
func (b *Board) SamePosition(that *Board) bool {
	for i := range b.brd {
		if b.brd[i] != that.brd[i] {
			return false
		}
	}
	return b.blk == that.blk && b.ep == that.ep && b.cas == that.cas
}

// Piece обозначает фигуру в клетке доски.
type Piece int

// Фигуры.
const (
	WhitePawn Piece = iota + 1
	BlackPawn
	WhiteKnight
	BlackKnight
	WhiteBishop
	BlackBishop
	WhiteRook
	BlackRook
	WhiteQueen
	BlackQueen
	WhiteKing
	BlackKing
)

// String возвращает строковое значение для фигуры, заглавные
// для белых и строчные для чёрных.
func (p Piece) String() string {
	switch p {
	case WhitePawn:
		return "P"
	case BlackPawn:
		return "p"
	case WhiteKnight:
		return "N"
	case BlackKnight:
		return "n"
	case WhiteBishop:
		return "B"
	case BlackBishop:
		return "b"
	case WhiteRook:
		return "R"
	case BlackRook:
		return "r"
	case WhiteQueen:
		return "Q"
	case BlackQueen:
		return "q"
	case WhiteKing:
		return "K"
	case BlackKing:
		return "k"
	default:
		return "-"
	}
}

// Castling обозначает рокировку или её возможность.
type Castling byte

const (
	WhiteKingside Castling = 1 << iota
	WhiteQueenside
	BlackKingside
	BlackQueenside
)

// Square - представление для клетки на доске.
type Square int8

var errSquareNotExist = fmt.Errorf("square does not exist")

// IsValid сообщает, есть ли на доске такая клетка.
func (s Square) IsValid() bool {
	return s >= 0 && s <= 63
}

// IsBlack возвращает true, если поле чёрное.
func (s Square) IsBlack() bool {
	return (s/8+s%8)%2 == 0
}

// String возвращает строковое представление клетки.
func (s Square) String() string {
	if s < 0 || s > 63 {
		return "-"
	}
	r := s / 8
	c := s % 8
	b := []byte{byte(c) + 'a', byte(r) + '1'}
	return string(b)
}

// Sq возвращает номер клетки на доске, -1 в случае, если клетки
// не существует. Принимает int от 0 до 63, либо строка вида
// "a1", "c7", "e8" и т.п.
// Соответственно, "a1" будет указывать на ту же клетку, что 0,
// "e1" - на ту же, что 4, и т. д., вплоть до "h8" - 63.
func Sq[S int | string](s S) Square {
	if ss, ok := any(s).(string); ok {
		if len(ss) != 2 {
			return -1
		}
		if ss[0] < 'a' || ss[0] > 'h' {
			return -1
		}
		if ss[1] < '1' || ss[1] > '8' {
			return -1
		}
		return Square(int(ss[1]-'1')*8 + int(ss[0]-'a'))
	}
	i := any(s).(int)
	if i < 0 || i > 63 {
		return -1
	}
	return Square(i)
}
