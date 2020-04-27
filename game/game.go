package game

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

const (
	SquareSize   float64 = 30
	BoardOffsetX float64 = 30
	BoardOffsetY float64 = 30
	ColCount     int     = 10
	RowCount     int     = 24
)

type Game struct {
	FilledArea []BoardPos
	ShapeImgs  *imdraw.IMDraw
	CurShape   *Shape

	NextShape     *Shape
	NextShapeImgs *imdraw.IMDraw

	IsGameOver bool
}

func New() *Game {
	board := Game{
		FilledArea:    make([]BoardPos, 0),
		ShapeImgs:     imdraw.New(nil),
		CurShape:      nil,
		NextShape:     NewShape(),
		NextShapeImgs: imdraw.New(nil),
		IsGameOver:    false,
	}

	board.SetShape(NewShape())
	return &board
}

func (game *Game) MoveShapeLeft() {
	for _, vtx := range game.CurShape.Coords {
		if vtx.X == 0 {
			return
		}

		for _, areaVtx := range game.FilledArea {
			if vtx.Y == areaVtx.Y && vtx.X-1 == areaVtx.X {
				return
			}
		}
	}

	for i := range game.CurShape.Coords {
		game.CurShape.Coords[i].X--
	}
}

func (game *Game) MoveShapeRight() {
	for _, vtx := range game.CurShape.Coords {
		if vtx.X == ColCount-1 {
			return
		}

		for _, areaVtx := range game.FilledArea {
			if vtx.Y == areaVtx.Y && vtx.X+1 == areaVtx.X {
				return
			}
		}
	}

	for i := range game.CurShape.Coords {
		game.CurShape.Coords[i].X++
	}
}

// Recalculate images positions from board state
func (game *Game) Refresh() {

	game.ShapeImgs.Clear()
	game.ShapeImgs.Color = colornames.Yellow

	for _, dot := range game.CurShape.Coords {
		game.ShapeImgs.Push(getSquarePos(dot)...)
		game.ShapeImgs.Rectangle(0)
	}

	game.ShapeImgs.Color = colornames.Green
	for _, dot := range game.FilledArea {
		game.ShapeImgs.Push(getSquarePos(dot)...)
		game.ShapeImgs.Rectangle(0)
	}

	game.NextShapeImgs.Clear()
	game.NextShapeImgs.Color = colornames.Yellow

	offsetVec := pixel.V(460, 750)
	for _, dot := range game.NextShape.Coords {
		pos := getSquarePos(dot)

		for idx, v := range pos {
			pos[idx] = v.To(offsetVec)
		}

		game.NextShapeImgs.Push(pos...)
		game.NextShapeImgs.Rectangle(0)
	}
}

func getSquarePos(pos BoardPos) (coords []pixel.Vec) {
	lB := pixel.V(BoardOffsetX+float64(pos.X)*SquareSize, BoardOffsetY+float64(pos.Y)*SquareSize)
	rB := pixel.V(lB.X+SquareSize, lB.Y)
	lT := pixel.V(lB.X, lB.Y+SquareSize)
	rT := pixel.V(rB.X, lT.Y)

	coords = []pixel.Vec{lB, rB, lT, rT}
	return
}

// Add new shape to the top of the field
func (game *Game) SetShape(shape *Shape) {
	game.CurShape = shape.ToBoardCoords()
}

func (game *Game) isPosFree(pos BoardPos) bool {
	if pos.Y < 0 {
		return false
	}

	for _, dot := range game.FilledArea {
		if pos.X == dot.X && pos.Y == dot.Y {
			return false
		}
	}

	return true
}

func (game *Game) RotateShape() {
	pivotIdx := game.CurShape.PivotIndex

	if pivotIdx < 0 {
		return
	}

	pivotCoords := game.CurShape.Coords[pivotIdx]
	var newShapePosition []BoardPos

	for _, vertex := range game.CurShape.Coords {
		// relative vector (abs position - pivot)
		vectorRelative := BoardPos{
			X: vertex.X - pivotCoords.X,
			Y: vertex.Y - pivotCoords.Y,
		}

		// transformation vr * rotation matrix
		vectorTransform := BoardPos{
			X: 0*vectorRelative.X + -1*vectorRelative.Y,
			Y: 1*vectorRelative.X + 0*vectorRelative.Y,
		}

		// transform + pivot
		newCoords := BoardPos{
			X: vectorTransform.X + pivotCoords.X,
			Y: vectorTransform.Y + pivotCoords.Y,
		}

		if !game.isPosFree(newCoords) {
			return
		}

		newShapePosition = append(newShapePosition, newCoords)
	}

	game.CurShape.Coords = newShapePosition

	// verify all vertices of figure are inside and move until shape is inside
	for _, vertex := range game.CurShape.Coords {
		if vertex.X < 0 {
			for i := vertex.X; i < 0; i++ {
				game.MoveShapeRight()
			}
		}

		if vertex.X > ColCount-1 {
			for i := vertex.X; i > ColCount-1; i-- {
				game.MoveShapeLeft()
			}
		}

		if vertex.Y < 0 {
			for i := vertex.Y; i == 0; i++ {

			}
		}
	}
}

func (game *Game) isColliding() bool {
	for _, figureVertex := range game.CurShape.Coords {
		if figureVertex.Y == 0 {
			return true // hit ground
		}

		for _, areaVertex := range game.FilledArea {
			if figureVertex.Y-1 == areaVertex.Y && figureVertex.X == areaVertex.X {
				return true // hit other piece
			}
		}
	}

	return false
}

func (game *Game) Fall() {
	if game.IsGameOver {
		return
	}

	if game.isColliding() {
		game.FilledArea = append(game.FilledArea, game.CurShape.Coords...)
		game.clearLines()
		game.SetShape(game.NextShape)
		game.NextShape = NewShape()
		return
	}

	for index := range game.CurShape.Coords {
		game.CurShape.Coords[index].Y--
	}
}

func fillGroup(i int, j int, groupNum int, board *[RowCount][ColCount]int) {
	if board[i][j] != -1 {
		return
	}
	board[i][j] = groupNum

	//top
	if i < RowCount-1 && board[i+1][j] == -1 {
		fillGroup(i+1, j, groupNum, board)
	}

	// bottom
	if i > 0 && board[i-1][j] == -1 {
		fillGroup(i-1, j, groupNum, board)
	}

	// left
	if j > 0 && board[i][j-1] == -1 {
		fillGroup(i, j-1, groupNum, board)
	}

	// right
	if j < ColCount-1 && board[i][j+1] == -1 {
		fillGroup(i, j+1, groupNum, board)
	}
}

func groupCanGoDown(board *[RowCount][ColCount]int, groupNum int) (canGoDown bool, groupDots []BoardPos) {
	for i := 0; i < RowCount; i++ {
		for j := 0; j < ColCount; j++ {
			if board[i][j] == groupNum {
				if i == 0 || (board[i-1][j] != 0 && board[i-1][j] != groupNum) {
					canGoDown = false
					groupDots = nil
					return
				}

				groupDots = append(groupDots, BoardPos{j, i})
			}
		}
	}

	canGoDown = true
	return
}

// Clearing lines using Sticky Gravity Mode
// Details: http://tetris.wikia.com/wiki/Line_clear
func (game *Game) clearLines() {
	var board [RowCount][ColCount]int
	for _, dot := range game.FilledArea {
		if dot.Y == RowCount-1 {
			fmt.Println("game over")
			game.IsGameOver = true
			return
		}
		board[dot.Y][dot.X] = -1
	}

	// Remove rows which are filled
	for i := 0; i < RowCount; i++ {
		rowFilled := true
		for j := 0; j < ColCount; j++ {
			if board[i][j] != -1 {
				rowFilled = false
				break
			}
		}

		if rowFilled {
			for j := 0; j < ColCount; j++ {
				board[i][j] = 0
			}
		}
	}

	// find groups using flood fill algorithm
	groupNum := 0
	for i := 0; i < RowCount; i++ {
		for j := 0; j < ColCount; j++ {
			if board[i][j] != -1 {
				continue
			}
			groupNum++
			fillGroup(i, j, groupNum, &board)
		}
	}

	// go through each group and move each group down until any dot of each group reaches the bottom or collides with other group
	for gr := 1; gr <= groupNum; gr++ {
		canGoDown, groupDots := groupCanGoDown(&board, gr)

		for canGoDown {
			for _, dot := range groupDots {
				board[dot.Y][dot.X] = 0
				board[dot.Y-1][dot.X] = gr
			}
			canGoDown, groupDots = groupCanGoDown(&board, gr)
		}
	}

	var newFilledArea []BoardPos

	for i := 0; i < RowCount; i++ {
		for j := 0; j < ColCount; j++ {
			if board[i][j] != 0 {
				newFilledArea = append(newFilledArea, BoardPos{j, i})
			}
		}
	}

	game.FilledArea = newFilledArea
}
