package domain

import (
	"fmt"
)

type TorusMatrix struct {
	width  int
	height int
}

func NewTorusMatrix(width, height int) (*TorusMatrix, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width and height must be positive integers, got width=%d, height=%d", width, height)
	}

	return &TorusMatrix{
		width:  width,
		height: height,
	}, nil
}

func (tm *TorusMatrix) Dimensions() (width, height int) {
	return tm.width, tm.height
}

func (tm *TorusMatrix) TotalElements() int {
	return tm.width * tm.height
}

func (tm *TorusMatrix) IndexToCoordinates(index int) (row, col int, err error) {
	if index < 0 || index >= tm.TotalElements() {
		return 0, 0, fmt.Errorf("index %d is out of bounds for matrix %dx%d", index, tm.width, tm.height)
	}

	row = index / tm.width
	col = index % tm.width
	return row, col, nil
}

func (tm *TorusMatrix) CoordinatesToIndex(row, col int) int {
	wrappedRow := ((row % tm.height) + tm.height) % tm.height
	wrappedCol := ((col % tm.width) + tm.width) % tm.width

	return wrappedRow*tm.width + wrappedCol
}

func (tm *TorusMatrix) IsValidIndex(index int) bool {
	return index >= 0 && index < tm.TotalElements()
}
