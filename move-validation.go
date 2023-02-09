package chessBox

import "fmt"

// ValidateMoveService сервис проверки хода: можно ли фигуру, которая находится на доске b на клетке startCell передвинуть на
// клетку endCell. На вход также подается пользователь o, чья очередь ходить.
// Возвращает pawnTransformation (true в случае, если пешка трансформируется в другую фигуру, достигнув
// конца доски, или false если это не так) и err (ошибку в случае невозможности сделать ход или nil, если ход возможен).
func ValidateMoveService(b *Board, o Opponent, startCell, endCell Cell) (pawnTransformation bool, err error) {
	// Логика валидации хода пошагово.

	// 1. получаем фигуру, находящуюся в клетке startCell. Если в этой клетке фигуры нет, возвращаем ошибку
	figure := b.GetFigure(startCell)
	if figure == nil {
		return pawnTransformation, fmt.Errorf("move invalid: startCell %d doesn't have any figures", startCell)
	}

	// 2. проверяем, что фигура принадлежит той стороне, чья очередь хода
	isFigureRightColor := checkFigureColor(o, figure)
	if !isFigureRightColor {
		return pawnTransformation, fmt.Errorf("move invalid: startCell %d has figure of wrong color", startCell)
	}

	// 3. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	var cellsToBePassed []Cell
	cellsToBePassed, err = checkFigureMove(figure, startCell, endCell)
	if err != nil {
		return pawnTransformation, fmt.Errorf("move invalid: %w", err)
	}

	// 4. Проверяем, что по пути фигуры с клетки startCell до endCell (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	if len(cellsToBePassed) > 0 {
		err = checkCellsToBePassed(b, cellsToBePassed)
		if err != nil {
			return pawnTransformation, fmt.Errorf("move invalid: %w", err)
		}
	}

	// 5. Проверяем наличие и цвет фигур в клетке endCell.
	err = checkEndCell(b, o, endCell)
	if err != nil {
		return pawnTransformation, fmt.Errorf("move invalid: %w", err)
	}

	// 6. На текущем этапе ход возможен. Генерируем новое положение доски newBoard.
	newBoard := b.GenerateBoardAfterMove(startCell, endCell)

	// 7. Проверяем, что при новой позиции на доске не появился шах для собственного короля оппонента o.
	err = checkSelfCheck(newBoard, o)
	if err != nil {
		return pawnTransformation, fmt.Errorf("move invalid: %w", err)
	}

	// 8. В случае если ход делается королем, проверяем, что он не подступил вплотную к чужому королю - такой ход
	// будет запрещен.

	// 9. В случае если ход делается пешкой, проверяем, попала ли пешка на последнюю линию и потребуется ли
	// трансформация в другую фигуру. Если да, выставляем pawnTransformation = true

	return pawnTransformation, nil
}

// isItOpponentTurn проверяет, что очередь хода оппонента o и цвет фигуры f, которую хотят передвинуть, совпадают.
// Возвращает true в случае успеха, false в противном случае.
func checkFigureColor(o Opponent, f Figure) bool {

}

// checkFigureMove проверяет, что фигура f может двигаться в этом направлении с клетки startCell на клетку
// endCell. Возвращает ошибку, если движение невозможно, или nil в противном случае. Также возвращают массив
// cellsToBePassed []Cell, в который входят все "промежуточные" клетки, которые "проходит" эта фигура на своем
// пути с startCell до endCell, если такие есть. В противном случае возвращает nil.
func checkFigureMove(f Figure, startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// через цепочку if идет проверка на тип Figure (или на поле name в структуре Figure) и перенаправление на
	// соответствующий метод Move для этого типа фигуры
}

// checkCellsToBePassed проверяет, есть ли на клетках из массива cellsToBePassed какие-либо фигуры. Если хотя бы на одной
// клетке есть фигура, возвращается ошибка. Иначе возвращается nil.
func checkCellsToBePassed(b *Board, cellsToBePassed []Cell) error {
	for _, cell := range cellsToBePassed {
		f := b.GetFigure(cell)
		if f != nil {
			return fmt.Errorf("figure present on cell %d", cell)
		}
	}
	return nil
}

// checkEndCell проверяет наличие фигуры на клетке endCell на предмет совместимости хода. Возвращает ошибку при
// несовместимости хода или nil в случае успеха.
func checkEndCell(b *Board, o Opponent, endCell Cell) error {
	// логика пошагово:

	// 1. Если в клетке endCell нет фигур, ход возможен:
	f := b.GetFigure(endCell)
	if f == nil {
		return nil
	}

	// 2. Если фигура в endCell принадлежит самому участнику o, ход невозможен.
	//

	// 3. Если фигура в endCell принадлежит сопернику, проверка, возможно ли взятие
	// (f.e. пешка e2-e4 не может взять коня на e4).
	//
}

// checkSelfCheck проверяет, что при новой позиции на доске нет шаха собственному королю.
func checkSelfCheck(b Board, o Opponent) error {
	// Находим клетку с собственным королем.
	kingCell := getKingCell(b, o)

	// проверяем, есть ли на расстоянии буквы Г от этой клетки вражеские кони. Если да, выдаем ошибку.
	err := checkEnemyKnightsNearKing(b, kingCell, o)
	if err != nil {
		return fmt.Errorf("self-check by enemy knight")
	}

	// проверяем, есть ли по вертикали или горизонтали в качестве ближайших фигур вражеские ладьи и ферзи.
	// Если да, выдаем ошибку.
	err = checkEnemiesVerticallyAndHorizontally(b, kingCell, o)
	if err != nil {
		return fmt.Errorf("self-check by enemy rook or queen")
	}

	// проверяем, есть ли по диагоналям в качестве ближайших фигур вражеские ферзи, слоны и пешки.
	// Если да, выдаем ошибку.
	err = checkEnemiesDiagonally(b, kingCell, o)
	if err != nil {
		return fmt.Errorf("self-check by enemy queen or pawn or bishop")
	}

	return nil

}

// getKingCell возвращает клетку, на которой находится свой король оппонента.
func getKingCell(b Board, o Opponent) Cell {

}

// checkEnemyKnightsNearKing проверяет ближайшие клетки в расположении буквой Г к своему королю на наличие на них
// вражеского коня. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemyKnightsNearKing(b Board, kingCell Cell, o Opponent) error {

}

// checkEnemiesVerticallyAndHorizontally проверяет ближайшие клетки по вертикали (сверху, снизу) и
// горизонтали (слева, справа) по отношению к клетке своего короля kingCell, на которых находятся вражеские ладьи и ферзи.
// Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesVerticallyAndHorizontally(b Board, kingCell Cell, o Opponent) error {

}

// checkEnemiesDiagonally проверяет ближайшие клетки по всем диагоналям по отношению к клетке своего короля kingCell, на
// которых находятся вражеские слоны, ферзи и пешки. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesDiagonally(b Board, kingCell Cell, o Opponent) error {
	// должна быть реализована дополнтельная проверка на пешки - их битое поле только по ходу движения!
}

// Методы Move для каждого типа фигуры. В случае невозможности сделать ход, возвращают ошибку, иначе возвращают nil.
// Также методы возвращают массив cellsToBePassed []Cell, в который входят все "промежуточные" клетки, которые "проходит"
// эта фигура на своем пути с startCell до endCell, если такие есть.

// Move для пешки
func (p *Pawn) Move(startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// логика движения пешки. Может двигаться вверх (или вниз в зависимости от цвета) на 1 или 2 клетки. Может съедать
	// по диагонали проходные пешки или фигуры соперника.
	// TODO: ПОДУМАТЬ КАК РЕАЛИЗОВАТЬ СЪЕДАНИЕ ПЕШКОЙ ПО ДИАГОНАЛИ, КОГДА ЧУЖАЯ ПЕШКА ИДЕТ ЧЕРЕЗ БИТОЕ ПОЛЕ! ДЛЯ ЭТОГО НУЖНО ИМЕТЬ ИНФОРМАЦИЮ О ПРОШЛОМ ХОДЕ!
}

// Move для короля
func (k *King) Move(startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// логика движения короля. Может двигаться вверх, вверх-вправо по диагонали, вправо, вправо-вниз по диагонали,
	// вниз, вниз-влево по диагонали, влево по диагонали, влево-вверх по диагонали. Только на одну клетку.

	// Также король может двигаться:
	// БЕЛЫЙ: на 2 клетки вправо (короткая рокировка) или три клетки влево (длинная рокировка)
	// ЧЕРНЫЙ: наоборот, на 2 клетки влево (короткая рокировка) или три клетки вправо (длинная рокировка). НО!
	// ТОЛЬКО ПРИ УСЛОВИИ, что
	// TODO: у структуры Король и ладья должен быть параметр, которые показывает, делали ли они уже ходы. Если да - рокировка невозможна!
	// TODO: для рокировки проходные поля для короля НЕ ДОЛЖНЫ быть под боем! там не должно быть под атакой исходное для короля, конечное для короля и то, через которое прыгает.
	// TODO: И ещё проверка на то что все поля между королем и ладьей свободны
	// TODO: коммент Алкесандра 2. то же самое и для пешки (чтобы делать первый ход на 2 клетки)
}

// Move для ферзя
func (q *Queen) Move(startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// логика движения королевы. Может двигаться диагонально на любое кол-во клеток. Мжно обеъдинить методы ладьи +
	// слона для ферзя. Либо можно ничего не объединять и прописать индивидуально для каждой фигуры
	// (возможны повторяющиеся куски кода).
}

// Move для ладьи
func (r *Rook) Move(startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// логика движения ладьи. Может двигаться вверх, вниз, влево, вправо на любое кол-во клеток.
}

// Move для слона
func (b *Bishop) Move(startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// логика движения слона. Может двигаться по всем диагоналям на любое кол-во клеток.
}

// Move для коня
func (k *Knight) Move(startCell, endCell Cell) (cellsToBePassed []Cell, err error) {
	// логика движения коня. Может двигаться буквой Г. То есть +/- 2 клетки в одном направлении и +/- 1 клетка
	// в перпендикулярном направении.
}

// Вспомогательные структуры

// Opponent моделирует пользователя с ником и цветом фигур, за которые он играет
type Opponent struct {
	name      string
	colorCode int
}

// Cell кетка доски
type Cell struct {
	row    int
	column int
}

// Board доска
type Board struct {
	Cells map[Cell]Figure
}

func (b *Board) NewBoard() {
	// NewBoard метод-конструктор доски в начале игры, создает доску и расставляет фигуры по правилам.
}

func (b *Board) Clear() {
	// удаляет все фигуры с доски
}

func (b *Board) GetFigure(c Cell) Figure {
	// возвращает фигуру, находщуюся на клетке c доски b. Если клетка пуста, возвращает nil.
}

func (b *Board) GenerateBoardAfterMove(startCell, endCell Cell) (newBoard Board) {
	// возвращает новую доску после движения фигуры с клетки startCell на клетку endCell.
	// В обычном случае, удаляет фигуру с клетки startCell и ставит эту же фигуру на клетку endCell.
	// В случае рокировки, отдельная логика для короткой и длинной.
}

// Figure интерфейс фигур с методом Move, уникальным для каждой фигуры. В метод передается значение начальной клетки
// startCell и клетки endCell, на которую планируется поставить фигуру. В случае невозможности это сделать,
// возвращается ошибка. При успехе возвращется nil.
type Figure interface {
	Move(startCell, endCell Cell) error
}

// структуры фигур (реализованы через интерфейс Figure)

type Pawn struct {
	colorCode int
	points    int
}

type King struct {
	colorCode int
	points    int
}

type Queen struct {
	colorCode int
	points    int
}

type Rook struct {
	colorCode int
	points    int
}

type Bishop struct {
	colorCode int
	points    int
}
type Knight struct {
	colorCode int
	points    int
}
