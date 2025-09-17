package domain

import "fmt"

type NeighborDirection struct {
	RowOffset int
	ColOffset int
	Name      string
}

var AllDirections = []NeighborDirection{
	{-1, -1, "TopLeft"},
	{-1, 0, "Top"},
	{-1, 1, "TopRight"},
	{0, -1, "Left"},
	{0, 1, "Right"},
	{1, -1, "BottomLeft"},
	{1, 0, "Bottom"},
	{1, 1, "BottomRight"},
}

type NeighborFinder struct {
	matrix *TorusMatrix
}

func NewNeighborFinder(matrix *TorusMatrix) *NeighborFinder {
	return &NeighborFinder{
		matrix: matrix,
	}
}

func (nf *NeighborFinder) FindNeighbors(index int) ([]int, error) {
	if !nf.matrix.IsValidIndex(index) {
		return nil, fmt.Errorf("invalid index %d for matrix dimensions %dx%d",
			index, nf.matrix.width, nf.matrix.height)
	}

	centerRow, centerCol, err := nf.matrix.IndexToCoordinates(index)
	if err != nil {
		return nil, fmt.Errorf("failed to convert index to coordinates: %w", err)
	}

	neighbors := make([]int, 0, len(AllDirections))

	for _, direction := range AllDirections {
		neighborRow := centerRow + direction.RowOffset
		neighborCol := centerCol + direction.ColOffset

		neighborIndex := nf.matrix.CoordinatesToIndex(neighborRow, neighborCol)
		neighbors = append(neighbors, neighborIndex)
	}

	return neighbors, nil
}
