package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"time"
)

const SquareSize float64 = 30
const BoardOffsetX float64 = 30
const BoardOffsetY float64 = 30
const ColCount int = 10
const RowCount int = 24

const ButtonRepeatTimeDelta float64 = 0.1

func mojaveWorkaround(win *pixelgl.Window) {
	// issue: https://github.com/faiface/pixel/issues/140
	pos := win.GetPos()
	win.SetPos(pixel.ZV)
	win.SetPos(pos)
}

func createBoardBgr(boardStartPos pixel.Vec, colCount int, rowCount int) *imdraw.IMDraw {
	boardGrid := imdraw.New(nil)
	boardGrid.Color = colornames.Lightgray

	for i := 0; i <= colCount; i++ {
		lineStart := pixel.V(float64(i)*SquareSize+boardStartPos.X, boardStartPos.Y)
		lineEnd := lineStart
		lineEnd.Y = boardStartPos.X + SquareSize*float64(rowCount)
		boardGrid.Push(lineStart, lineEnd)
		boardGrid.Line(1.1)
	}

	for j := 0; j <= rowCount; j++ {
		lineStart := pixel.V(boardStartPos.X, boardStartPos.Y+float64(j)*SquareSize)
		lineEnd := lineStart
		lineEnd.X = boardStartPos.X + float64(colCount)*SquareSize
		boardGrid.Push(lineStart, lineEnd)
		boardGrid.Line(1.1)
	}

	return boardGrid
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GoTetris",
		Bounds: pixel.R(0, 0, 500, 800),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	mojaveWorkaround(window)
	window.SetSmooth(false)

	boardBgr := createBoardBgr(pixel.V(BoardOffsetX, BoardOffsetY), ColCount, RowCount)
	gameBoard := NewGameBoard()

	last := time.Now()
	lastPressed := last
	pressedDt := time.Since(lastPressed).Seconds()

	lastFall := time.Now()
	fallDt := time.Since(lastFall).Seconds()

	fallSpeed := 0.2

	for !window.Closed() {
		if gameBoard.IsGameOver {
			window.SetTitle("Game Over!")
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		if window.Pressed(pixelgl.KeyLeft) {
			pressedDt = time.Since(lastPressed).Seconds()
			if pressedDt >= ButtonRepeatTimeDelta {
				gameBoard.MoveShapeLeft()
				lastPressed = time.Now()
			}

		}

		if window.Pressed(pixelgl.KeyRight) {
			pressedDt = time.Since(lastPressed).Seconds()
			if pressedDt >= ButtonRepeatTimeDelta {
				gameBoard.MoveShapeRight()
				lastPressed = time.Now()
			}
		}

		if window.JustPressed(pixelgl.KeyUp) {
			gameBoard.RotateShape()
		}

		if window.Pressed(pixelgl.KeyDown) {
			fallSpeed = 0.05
		} else {
			fallSpeed = 0.2
		}

		fallDt = time.Since(lastFall).Seconds()
		if fallDt >= fallSpeed {
			lastFall = time.Now()
			gameBoard.fallDown()
		}

		if dt >= 1/60 {
			window.Clear(colornames.Black)
			boardBgr.Draw(window)
			gameBoard.Refresh()
			gameBoard.ShapeImages.Draw(window)
			gameBoard.NextShapeImages.Draw(window)
			window.Update()
		}
	}
}

func main() {
	pixelgl.Run(run)
}
