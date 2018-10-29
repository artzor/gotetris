package main

import (
	"math/rand"
	"time"
)

type BoardPos struct {
	X int
	Y int
}

type Shape struct {
	PivotIndex int
	Coords     []BoardPos
}

func NewShape() (shape *Shape) {
	shape = &Shape{}
	pos, pivotIndex := getRandomShape()

	shape.PivotIndex = pivotIndex
	shape.Coords = pos
	return
}

func getRandomShape() (pos []BoardPos, pivotIndex int) {
	rand.Seed(time.Now().UTC().UnixNano())
	shapes := []func() (pos []BoardPos, pivotIndex int){
		getShapeL,
		getShapeJ,
		getShapeT,
		getShapeI,
		getShapeZ,
		getShapeS,
		getShapeO,
	}

	idx := rand.Intn(len(shapes))
	shapeFunc := shapes[idx]
	pos, pivotIndex = shapeFunc()
	return
}

func getShapeL() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 0},
		{0, 1},
		{0, 2},
		{1, 0},
	}

	pivotIndex = 1
	return
}

func getShapeJ() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 0},
		{1, 0},
		{1, 1},
		{1, 2},
	}

	pivotIndex = 2
	return
}

func getShapeT() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 0},
		{0, 1},
		{0, 2},
		{1, 1},
	}

	pivotIndex = 1
	return
}

func getShapeI() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	}

	pivotIndex = 2
	return
}

func getShapeZ() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 1},
		{1, 0},
		{1, 1},
		{2, 0},
	}

	pivotIndex = 2
	return
}

func getShapeS() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	}

	pivotIndex = 2
	return
}

func getShapeO() (pos []BoardPos, pivotIndex int) {
	pos = []BoardPos{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
	}

	pivotIndex = -1
	return
}
