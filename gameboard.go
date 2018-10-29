package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type GameBoard struct {
	FilledArea   []BoardPos
	ShapeImages  *imdraw.IMDraw
	CurrentShape *Shape

	NextShape       *Shape
	NextShapeImages *imdraw.IMDraw

	IsGameOver bool
}

func NewGameBoard() *GameBoard {
	board := GameBoard{}

	board.ShapeImages = imdraw.New(nil)
	board.FilledArea = make([]BoardPos, 0)
	board.setCurrentShape(NewShape())

	board.NextShape = NewShape()
	board.NextShapeImages = imdraw.New(nil)
	board.IsGameOver = false

	return &board
}

func (gameBoard *GameBoard) MoveShapeLeft() {
	for _, vertex := range gameBoard.CurrentShape.Coords {
		if vertex.X == 0 {
			return
		}

		for _, areaVertex := range gameBoard.FilledArea {
			if vertex.Y == areaVertex.Y && vertex.X-1 == areaVertex.X {
				return
			}
		}
	}

	for index := range gameBoard.CurrentShape.Coords {
		gameBoard.CurrentShape.Coords[index].X--
	}
}

func (gameBoard *GameBoard) MoveShapeRight() {
	for _, vertex := range gameBoard.CurrentShape.Coords {
		if vertex.X == ColCount-1 {
			return
		}

		for _, areaVertex := range gameBoard.FilledArea {
			if vertex.Y == areaVertex.Y && vertex.X+1 == areaVertex.X {
				return
			}
		}
	}

	for index := range gameBoard.CurrentShape.Coords {
		gameBoard.CurrentShape.Coords[index].X++
	}
}

// Recalculate images positions from board state
func (gameBoard *GameBoard) Refresh() {

	gameBoard.ShapeImages.Clear()
	gameBoard.ShapeImages.Color = colornames.Yellow

	for _, dot := range gameBoard.CurrentShape.Coords {
		gameBoard.ShapeImages.Push(getSquarePos(dot)...)
		gameBoard.ShapeImages.Rectangle(0)
	}

	gameBoard.ShapeImages.Color = colornames.Green
	for _, dot := range gameBoard.FilledArea {
		gameBoard.ShapeImages.Push(getSquarePos(dot)...)
		gameBoard.ShapeImages.Rectangle(0)
	}

	gameBoard.NextShapeImages.Clear()
	gameBoard.NextShapeImages.Color = colornames.Yellow

	offsetVec := pixel.V(460, 750)
	for _, dot := range gameBoard.NextShape.Coords {
		pos := getSquarePos(dot)

		for idx, v := range pos {
			pos[idx] = v.To(offsetVec)
		}

		gameBoard.NextShapeImages.Push(pos...)
		gameBoard.NextShapeImages.Rectangle(0)
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
func (gameBoard *GameBoard) setCurrentShape(shape *Shape) {
	row := RowCount - 3
	col := ColCount/2 - 1

	for index := range shape.Coords {
		shape.Coords[index].X = shape.Coords[index].X + col
		shape.Coords[index].Y = shape.Coords[index].Y + row
	}

	gameBoard.CurrentShape = shape
}

func (gameBoard *GameBoard) isPosFree(pos BoardPos) bool {
	if pos.Y < 0 {
		return false
	}

	for _, dot := range gameBoard.FilledArea {
		if pos.X == dot.X && pos.Y == dot.Y {
			return false
		}
	}

	return true
}

func (gameBoard *GameBoard) RotateShape() {
	pivotIdx := gameBoard.CurrentShape.PivotIndex

	if pivotIdx < 0 {
		return
	}

	pivotCoords := gameBoard.CurrentShape.Coords[pivotIdx]
	var newShapePosition []BoardPos

	for _, vertex := range gameBoard.CurrentShape.Coords {
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

		if !gameBoard.isPosFree(newCoords) {
			return
		}

		newShapePosition = append(newShapePosition, newCoords)
	}

	gameBoard.CurrentShape.Coords = newShapePosition

	// verify all vertexes of figure are inside and move until shape is inside
	for _, vertex := range gameBoard.CurrentShape.Coords {
		if vertex.X < 0 {
			for i := vertex.X; i < 0; i++ {
				gameBoard.MoveShapeRight()
			}
		}

		if vertex.X > ColCount-1 {
			for i := vertex.X; i > ColCount-1; i-- {
				gameBoard.MoveShapeLeft()
			}
		}

		if vertex.Y < 0 {
			for i := vertex.Y; i == 0; i++ {

			}
		}
	}
}

func (gameBoard *GameBoard) isColliding() bool {
	for _, figureVertex := range gameBoard.CurrentShape.Coords {
		if figureVertex.Y == 0 {
			return true // hit ground
		}

		for _, areaVertex := range gameBoard.FilledArea {
			if figureVertex.Y-1 == areaVertex.Y && figureVertex.X == areaVertex.X {
				return true // hit other piece
			}
		}
	}

	return false
}

func (gameBoard *GameBoard) fallDown() {
	if gameBoard.IsGameOver {
		return
	}

	if gameBoard.isColliding() {
		gameBoard.FilledArea = append(gameBoard.FilledArea, gameBoard.CurrentShape.Coords...)
		gameBoard.clearLines()
		gameBoard.setCurrentShape(gameBoard.NextShape)
		gameBoard.NextShape = NewShape()
		return
	}

	for index := range gameBoard.CurrentShape.Coords {
		gameBoard.CurrentShape.Coords[index].Y--
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
func (gameBoard *GameBoard) clearLines() {
	var board [RowCount][ColCount]int
	for _, dot := range gameBoard.FilledArea {
		if dot.Y == RowCount-1 {
			fmt.Println("game over")
			gameBoard.IsGameOver = true
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

	gameBoard.FilledArea = newFilledArea
}
