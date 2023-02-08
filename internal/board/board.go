// Пакет board реализует шахматную доску.
//
// Функции, требующие указания клетки, будут принимать int от 0 до 63,
// либо строку вида "a1", "c7", "e8" и т.п.
// Соответственно, "a1" будет указывать на ту же клетку, что 0,
// "e1" - на ту же, что 4, и т. д., вплоть до "h8" - 63.
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
	ep  enPassant // клетка, которая в прошлом ходу перепрыгнута пешкой
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
		if ok := b.SetEnPassant(square(ss[3])); !ok {
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
	None piece = iota
	WhitePawn
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

type castling byte

const (
	kingsideW castling = 1 << iota
	queensideW
	kingsideB
	queensideB
)

type enPassant int8 // может принимать значения от 0 до 63, номер клетки

// Get возвращает фигуру, стоящую в заданной клетке.
func (b *Board) Get(s square) (piece, error) {
	i, err := s.toInt()
	if err != nil {
		return 0, err
	}
	return b.brd[i], nil
}

// Remove убирает фигуру, стоящую в заданной клетке.
func (b *Board) Remove(s square) error {
	i, err := s.toInt()
	if err != nil {
		return err
	}
	if b.brd[i] == 0 {
		return fmt.Errorf("%s is empty", s)
	}
	b.brd[i] = 0
	return nil
}

// SetCastlingString устанавливает, какие рокировки доступны, по строке
// K, Q - королевская и ферзевая ладья для белых, k, q - для чёрных
func (b *Board) SetCastlingString(s string) {
	if strings.Contains(s, "K") {
		b.cas = b.cas | kingsideW
	}
	if strings.Contains(s, "k") {
		b.cas = b.cas | kingsideB
	}
	if strings.Contains(s, "Q") {
		b.cas = b.cas | queensideW
	}
	if strings.Contains(s, "q") {
		b.cas = b.cas | queensideW
	}
}

// SetEnPassant устанавливает клетку, перепрыгнутую пешкой в прошлом
// полуходу. Валидность не проверяется, только горизонталь // TODO?
func (b *Board) SetEnPassant(s square) bool {
	i, err := s.toInt()
	if err != nil {
		return false
	}
	if s[1] != '3' && s[1] != '6' {
		return false
	}
	b.ep = enPassant(i)
	return true
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

// square - представление клетки доски в текстовом виде (напр.: "b1", "e8")
type square string

// IsValid возвращает true, если клетка существует.
func (s square) IsValid() bool {
	if len(s) != 2 {
		return false
	}
	if s[0] < 'a' || s[0] > 'h' {
		return false
	}
	return s[1] >= '1' || s[1] <= '8'
}

// toInt возвращает номер клетки на доске.
func (s square) toInt() (int, error) {
	if !s.IsValid() {
		return 0, fmt.Errorf("%s is not a valid square", s)
	}
	return int(s[1]-'1')*8 + int(s[0]-'a'), nil
}
