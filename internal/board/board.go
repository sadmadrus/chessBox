// Пакет board реализует шахматную доску.
package board

import (
	"fmt"
	"strconv"
	"strings"
)

// Board представляет доску с позицией. Инициализация необязательна,
// пустой Board готов к употреблению, как пустая доска (без фигур,
// следующий ход первый).
type Board struct {
	brd [64]piece // доска из 64 клеток
	blk bool      // false - ход белых, true - чёрных
	cas castling  // битовая маска возможных рокировок
	ep  square    // клетка, которая в прошлом ходу перепрыгнута пешкой
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
	if ss[2] != "-" {
		if ok := b.SetEnPassant(Sq(ss[3])); !ok {
			return nil, er
		}
	}
	hm, err := strconv.Atoi(ss[4])
	if err != nil {
		return nil, er
	}
	fm, err := strconv.Atoi(ss[5])
	if err != nil {
		return nil, er
	}
	b.hm = hm
	b.fm = fm

	b.SetCastlingString(ss[2])

	ss = strings.Split(ss[0], "/")
	if len(ss) != 8 {
		return nil, er
	}
	for row := 0; row < 8; row++ {
		c := 0
		for _, r := range ss[7-row] {
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
				v := int(r - '0')
				if v < 1 || v > 8 {
					return nil, er
				}
				c += v - 1
			}
			c++
		}
	}
	return b, nil
}

// piece обозначает фигуру в клетке доски.
type piece int

const (
	WhitePawn piece = iota + 1
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

// castling обозначает возможность рокировки.
type castling byte

const (
	CastlingWhiteKingside castling = 1 << iota
	CastlingWhiteQueenside
	CastlingBlackKingside
	CastlingBlackQueenside
)

// Get возвращает фигуру, стоящую в заданной клетке.
func (b *Board) Get(s square) (piece, error) {
	if s == -1 {
		return 0, fmt.Errorf("%w: %v", errSquareNotExist, s)
	}
	return b.brd[s], nil
}

// Move перемещает фигуру из from в to.
// Валидность хода не проверяется, но по итогам ход переходит к
// другому игроку, количество полуходов и ходов обновляется,
// обновляются данные о пешке en passant и рокировках.
// Чтобы сделать рокировку, надо, например, сделать ход королём,
// а ладью убрать с доски и поставить на новое место, чтобы это
// защиталось как только один ход.
func (b *Board) Move(from, to square) error {
	if from < 0 {
		return fmt.Errorf("%w: %v", errSquareNotExist, from)
	}
	if to < 0 {
		return fmt.Errorf("%w: %v", errSquareNotExist, to)
	}
	b.ep = 0
	if b.blk {
		b.fm++
	}
	b.blk = !b.blk
	pc := b.brd[from]
	if pc != BlackPawn && pc != WhitePawn && b.brd[to] == 0 {
		b.hm++
	}
	switch from {
	case 0:
		b.RemoveCastling(CastlingWhiteQueenside)
	case 7:
		b.RemoveCastling(CastlingWhiteKingside)
	case 56:
		b.RemoveCastling(CastlingBlackQueenside)
	case 63:
		b.RemoveCastling(CastlingBlackKingside)
	}
	switch pc {
	case BlackKing:
		b.RemoveCastling(CastlingBlackKingside)
		b.RemoveCastling(CastlingBlackQueenside)
	case WhiteKing:
		b.RemoveCastling(CastlingWhiteKingside)
		b.RemoveCastling(CastlingWhiteQueenside)
	case BlackPawn:
		if from >= 8*6 && from < 8*7 && to >= 8*4 && to < 8*5 {
			b.ep = 8*5 + to%8
		}
	case WhitePawn:
		if from >= 8*1 && from < 8*2 && to >= 8*3 && to < 8*4 {
			b.ep = 8*2 + to%8
		}
	}
	b.brd[from] = 0
	b.brd[to] = pc
	return nil
}

// Remove убирает фигуру, стоящую в заданной клетке.
func (b *Board) Remove(s square) error {
	if s == -1 {
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

// BlackToMove устанавливает следующих ход белых.
func (b *Board) BlackToMove() {
	b.blk = true
}

// NextToMove возвращает true, если следующий ход белых,
// false, если следующий ход чёрных.
func (b *Board) NextToMove() bool {
	return !b.blk
}

// SetCastlingString устанавливает, какие рокировки доступны, по строке
// K, Q - королевская и ферзевая ладья для белых, k, q - для чёрных
// TODO: проверять, что в строке нет лишних символов
func (b *Board) SetCastlingString(s string) {
	if strings.Contains(s, "K") {
		b.SetCastling(CastlingWhiteKingside)
	}
	if strings.Contains(s, "k") {
		b.SetCastling(CastlingBlackKingside)
	}
	if strings.Contains(s, "Q") {
		b.SetCastling(CastlingWhiteQueenside)
	}
	if strings.Contains(s, "q") {
		b.SetCastling(CastlingBlackQueenside)
	}
}

// SetCastling устанавливает доступность рокировки.
func (b *Board) SetCastling(c castling) {
	b.cas = b.cas | c
}

// HaveCastling проверяет доступность рокировки.
func (b *Board) HaveCastling(c castling) bool {
	return b.cas&c != 0
}

// RemoveCastling убирает возможность рокировки.
func (b *Board) RemoveCastling(c castling) {
	b.cas = b.cas &^ c
}

// SetEnPassant устанавливает клетку, перепрыгнутую пешкой в прошлом
// полуходу. Валидность не проверяется, только горизонталь // TODO?
func (b *Board) SetEnPassant(s square) bool {
	if s/8 != 2 && s/8 != 5 {
		return false
	}
	b.ep = s
	return true
}

// GetEnPassant возвращает клетку, перепрыгнутую пешкой в прошлом ходу,
// -1 если такой не было.
func (b *Board) GetEnPassant() square {
	if b.ep == 0 {
		return -1
	}
	return b.ep
}

// IsEnPassant возвращает true, если заданная клетка была перепрыгнута
// пешкой в прошлом ходу.
func (b *Board) IsEnPassant(s square) bool {
	if s == 0 {
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

// square - представление для клетки на доске.
type square int8

var errSquareNotExist = fmt.Errorf("square does not exist")

// String возвращает строковое представление клетки.
func (s square) String() string {
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
func Sq[S int | string](s S) square {
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
		return square(int(ss[1]-'1')*8 + int(ss[0]-'a'))
	}
	i := any(s).(int)
	if i < 0 || i > 63 {
		return -1
	}
	return square(i)
}
