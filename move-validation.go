// пакет validation валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре.

package validation

import (
	"fmt"
	"net/http"
)

// http хендлеры

// SimpleValidation сервис отвечает за простую валидацию хода по начальной и конечной клетке и фигуре из URL-параметров запроса.
// Возвращает заголовок HttpResponse 200 (ход валиден) или HttpsResponse 403 (ход невалиден).
func SimpleValidation(w http.ResponseWriter, r *http.Request) {

}

// AdvancedValidation сервис отвечает за сложную валидацию хода по начальной и конечной клетке, а также по текущему состоянию
// доски в нотации FEN. Также принимает на вход URL-параметр newpiece (это новая фигура, в которую нужно превратить
// пешку при достижении последнего ряда, в формате pieceВозвращает заголовок HttpResponse 200 (ход валиден) или
// HttpsResponse 403 (ход невалиден). Возвращает в теле JSON с конечной доской board в форате FEN.
func AdvancedValidation(w http.ResponseWriter, r *http.Request) {

}

// GetAvailableMoves сервис отвечает за оплучение всех возможных ходов для данной позиции доски в нотации FEN и начальной клетке.
// Возвращает заголовок HttpResponse 200 (в случае непустого массива клеток) или HttpsResponse 403 (клетка пустая или
// с фигурой, которой не принадлежит ход или массив клеток пуст). Возвращает в теле
// JSON массив всех клеток, движение на которые валидно для данной фигуры.
func GetAvailableMoves(w http.ResponseWriter, r *http.Request) {

}

// Вспомогательные функции

// ValidateMove обрабывает общую логику валидации хода
func ValidateMove(b Board, from, to s) error {
	// Логика валидации хода пошагово.

	// 1. получаем фигуру, находящуюся в клетке startCell. Если в этой клетке фигуры нет, возвращаем ошибку
	figure := b.GetFigure(from)
	if figure == nil {
		return fmt.Errorf("move invalid: startCell %d doesn't have any figures", from)
	}

	// 2. проверяем, что фигура принадлежит той стороне, чья очередь хода
	isFigureRightColor := checkFigureColor(figure)
	if !isFigureRightColor {
		return fmt.Errorf("move invalid: startCell %d has figure of wrong color", from)
	}

	// 3. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	var cellsToBePassed []square
	cellsToBePassed, err := checkFigureMove(figure, from, to)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// 4. Проверяем, что по пути фигуры с клетки startCell до endCell (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	if len(cellsToBePassed) > 0 {
		err = checkCellsToBePassed(b, cellsToBePassed)
		if err != nil {
			return fmt.Errorf("move invalid: %w", err)
		}
	}

	// 5. Проверяем наличие и цвет фигур в клетке endCell.
	err = checkEndCell(b, to)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// TODO: проверяем, что пользователь не захотел выставить нового ферзя в неуместном для этого случае.

	// 6. На текущем этапе ход возможен. Генерируем новое положение доски newBoard.
	newBoard := b.GenerateBoardAfterMove(from, to)

	// 7. Проверяем, что при новой позиции на доске не появился шах для собственного короля оппонента o.
	err = checkSelfCheck(newBoard)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// 8. В случае если ход делается королем, проверяем, что он не подступил вплотную к чужому королю - такой ход
	// будет запрещен.

	// 9. В случае если ход делается пешкой, проверяем, попала ли пешка на последнюю линию и потребуется ли
	// трансформация в другую фигуру. Если да, выставляем pawnTransformation = true

	return nil
}

// checkFigureColor проверяет, что очередь хода и цвет фигуры f, которую хотят передвинуть, совпадают.
// Возвращает true в случае успеха, false в противном случае.
func checkFigureColor(b Board) bool {

}

// checkFigureMove проверяет, что фигура f может двигаться в этом направлении с клетки startCell на клетку
// endCell. Возвращает ошибку, если движение невозможно, или nil в противном случае. Также возвращают массив
// cellsToBePassed []Cell, в который входят все "промежуточные" клетки, которые "проходит" эта фигура на своем
// пути с startCell до endCell, если такие есть. В противном случае возвращает nil.
func checkFigureMove(f Figure, from, to s) (cellsToBePassed []Cell, err error) {
	// через цепочку if идет проверка на тип Figure (или на поле name в структуре Figure) и перенаправление на
	// соответствующий метод Move для этого типа фигуры
}

// checkCellsToBePassed проверяет, есть ли на клетках из массива cellsToBePassed какие-либо фигуры. Если хотя бы на одной
// клетке есть фигура, возвращается ошибка. Иначе возвращается nil.
func checkCellsToBePassed(b Board, cellsToBePassed []Cell) error {
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
func checkEndCell(b *Board, to s) error {
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
func checkSelfCheck(b Board) error {
	// Находим клетку с собственным королем.
	kingCell := getKingCell(b)

	// проверяем, есть ли на расстоянии буквы Г от этой клетки вражеские кони. Если да, выдаем ошибку.
	err := checkEnemyKnightsNearKing(b, kingCell)
	if err != nil {
		return fmt.Errorf("self-check by enemy knight")
	}

	// проверяем, есть ли по вертикали или горизонтали в качестве ближайших фигур вражеские ладьи и ферзи.
	// Если да, выдаем ошибку.
	err = checkEnemiesVerticallyAndHorizontally(b, kingCell)
	if err != nil {
		return fmt.Errorf("self-check by enemy rook or queen")
	}

	// проверяем, есть ли по диагоналям в качестве ближайших фигур вражеские ферзи, слоны и пешки.
	// Если да, выдаем ошибку.
	err = checkEnemiesDiagonally(b, kingCell)
	if err != nil {
		return fmt.Errorf("self-check by enemy queen or pawn or bishop")
	}

	return nil

}

// getKingCell возвращает клетку, на которой находится свой король оппонента.
func getKingCell(b Board) s {

}

// checkEnemyKnightsNearKing проверяет ближайшие клетки в расположении буквой Г к своему королю на наличие на них
// вражеского коня. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemyKnightsNearKing(b Board, kingCell Cell) error {

}

// checkEnemiesVerticallyAndHorizontally проверяет ближайшие клетки по вертикали (сверху, снизу) и
// горизонтали (слева, справа) по отношению к клетке своего короля kingCell, на которых находятся вражеские ладьи и ферзи.
// Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesVerticallyAndHorizontally(b Board, kingCell Cell) error {

}

// checkEnemiesDiagonally проверяет ближайшие клетки по всем диагоналям по отношению к клетке своего короля kingCell, на
// которых находятся вражеские слоны, ферзи и пешки. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesDiagonally(b Board, kingCell Cell) error {
	// должна быть реализована дополнтельная проверка на пешки - их битое поле только по ходу движения!
}

// Методы Move для каждого типа фигуры. В случае невозможности сделать ход, возвращают ошибку, иначе возвращают nil.

// MovePawn для пешки
func MovePawn(from, to s) error {
	// логика движения пешки. Может двигаться вверх (или вниз в зависимости от цвета) на 1 или 2 клетки. Может съедать
	// по диагонали проходные пешки или фигуры соперника.
}

// MoveKing для короля
func MoveKing(from, to s) error {
	// логика движения короля. Может двигаться вверх, вверх-вправо по диагонали, вправо, вправо-вниз по диагонали,
	// вниз, вниз-влево по диагонали, влево по диагонали, влево-вверх по диагонали. Только на одну клетку.

	// Также король может двигаться:
	// БЕЛЫЙ: на 2 клетки вправо (короткая рокировка) или 2 клетки влево (длинная рокировка)
	// ЧЕРНЫЙ: наоборот, на 2 клетки влево (короткая рокировка) или 2 клетки вправо (длинная рокировка). НО!
	// TODO: для рокировки проходные поля для короля НЕ ДОЛЖНЫ быть под боем! там не должно быть под атакой исходное для короля, конечное для короля и то, через которое прыгает.
	// TODO: И ещё проверка на то что все поля между королем и ладьей свободны
}

// MoveQueen для ферзя
func MoveQueen(from, to s) error {
	// логика движения королевы. Может двигаться диагонально на любое кол-во клеток. Мжно обеъдинить методы ладьи +
	// слона для ферзя. Либо можно ничего не объединять и прописать индивидуально для каждой фигуры
	// (возможны повторяющиеся куски кода).
}

// MoveRook для ладьи
func MoveRook(from, to s) error {
	// логика движения ладьи. Может двигаться вверх, вниз, влево, вправо на любое кол-во клеток.
}

// MoveBishop для слона
func MoveBishop(from, to s) error {
	// логика движения слона. Может двигаться по всем диагоналям на любое кол-во клеток.
}

// MoveKnight для коня
func MoveKnight(from, to s) error {
	// логика движения коня. Может двигаться буквой Г. То есть +/- 2 клетки в одном направлении и +/- 1 клетка
	// в перпендикулярном направении.
}
