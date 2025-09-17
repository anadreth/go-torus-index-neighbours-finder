package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type MatrixHasher struct {
	matrix *TorusMatrix
}

func NewMatrixHasher(matrix *TorusMatrix) *MatrixHasher {
	return &MatrixHasher{
		matrix: matrix,
	}
}

func (mh *MatrixHasher) GenerateWrappedMatrix() [][]int {
	width, height := mh.matrix.Dimensions()

	extendedHeight := height + 2
	extendedWidth := width + 2

	extended := make([][]int, extendedHeight)
	for i := range extended {
		extended[i] = make([]int, extendedWidth)
	}

	for extRow := 0; extRow < extendedHeight; extRow++ {
		for extCol := 0; extCol < extendedWidth; extCol++ {
			originalRow := extRow - 1
			originalCol := extCol - 1

			wrappedIndex := mh.matrix.CoordinatesToIndex(originalRow, originalCol)
			extended[extRow][extCol] = wrappedIndex
		}
	}

	return extended
}

func (mh *MatrixHasher) GenerateMatrixString() string {
	wrapped := mh.GenerateWrappedMatrix()

	var rows []string
	for _, row := range wrapped {
		var elements []string
		for _, element := range row {
			elements = append(elements, strconv.Itoa(element))
		}
		rows = append(rows, strings.Join(elements, ","))
	}

	return strings.Join(rows, "\n")
}

func (mh *MatrixHasher) CalculateHash() string {
	matrixString := mh.GenerateMatrixString()

	hasher := sha256.New()
	hasher.Write([]byte(matrixString))
	hashBytes := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashBytes)
}

func (mh *MatrixHasher) ValidateExpectedHash(expected string) error {
	calculated := mh.CalculateHash()
	if calculated != expected {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expected, calculated)
	}
	return nil
}
