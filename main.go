package main

import (
	"gotetris/game"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	slow = 0.2
	fast = 0.05
)

const ButtonRepeatTimeDelta float64 = 0.1

func macOSFix(win *pixelgl.Window) {
	// issue: https://github.com/faiface/pixel/issues/140
	pos := win.GetPos()
	win.SetPos(pixel.ZV)
	win.SetPos(pos)
}

func makeBgr(boardStartPos pixel.Vec, colCount int, rowCount int) *imdraw.IMDraw {
	boardGrid := imdraw.New(nil)
	boardGrid.Color = colornames.Lightgray

	for i := 0; i <= colCount; i++ {
		lineStart := pixel.V(float64(i)*game.SquareSize+boardStartPos.X, boardStartPos.Y)
		lineEnd := lineStart
		lineEnd.Y = boardStartPos.X + game.SquareSize*float64(rowCount)
		boardGrid.Push(lineStart, lineEnd)
		boardGrid.Line(1.1)
	}

	for j := 0; j <= rowCount; j++ {
		lineStart := pixel.V(boardStartPos.X, boardStartPos.Y+float64(j)*game.SquareSize)
		lineEnd := lineStart
		lineEnd.X = boardStartPos.X + float64(colCount)*game.SquareSize
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
		log.Panicf("failed to create window: %v", err)
	}
	macOSFix(window)
	window.SetSmooth(false)

	g := game.New()
	boardBgr := makeBgr(pixel.V(game.BoardOffsetX, game.BoardOffsetY), game.ColCount, game.RowCount)

	last := time.Now()
	lastPressed := last
	pressedDt := time.Since(lastPressed).Seconds()

	lastFall := time.Now()
	fallDt := time.Since(lastFall).Seconds()

	fallSpeed := 0.2

	for !window.Closed() {
		if g.IsGameOver {
			window.SetTitle("Game Over!")
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		if window.Pressed(pixelgl.KeyLeft) {
			pressedDt = time.Since(lastPressed).Seconds()
			if pressedDt >= ButtonRepeatTimeDelta {
				g.MoveShapeLeft()
				lastPressed = time.Now()
			}

		}

		if window.Pressed(pixelgl.KeyRight) {
			pressedDt = time.Since(lastPressed).Seconds()
			if pressedDt >= ButtonRepeatTimeDelta {
				g.MoveShapeRight()
				lastPressed = time.Now()
			}
		}

		if window.JustPressed(pixelgl.KeyUp) {
			g.RotateShape()
		}

		if window.Pressed(pixelgl.KeyDown) {
			fallSpeed = fast
		} else {
			fallSpeed = slow
		}

		fallDt = time.Since(lastFall).Seconds()
		if fallDt >= fallSpeed {
			lastFall = time.Now()
			g.Fall()
		}

		if dt >= 1/60 {
			window.Clear(colornames.Black)
			boardBgr.Draw(window)
			g.Refresh()
			g.ShapeImgs.Draw(window)
			g.NextShapeImgs.Draw(window)
			window.Update()
		}
	}
}

func main() {
	pixelgl.Run(run)
}
