package domain

import (
	"strings"
	"testing"
)

func TestGenerateWrappedMatrix(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected [][]int
	}{
		{
			name:   "4x4 matrix from problem example",
			width:  4,
			height: 4,
			expected: [][]int{
				{15, 12, 13, 14, 15, 12},
				{3, 0, 1, 2, 3, 0},
				{7, 4, 5, 6, 7, 4},
				{11, 8, 9, 10, 11, 8},
				{15, 12, 13, 14, 15, 12},
				{3, 0, 1, 2, 3, 0},
			},
		},
		{
			name:   "1x1 matrix",
			width:  1,
			height: 1,
			expected: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name:   "2x3 matrix",
			width:  3,
			height: 2,
			expected: [][]int{
				{5, 3, 4, 5, 3},
				{2, 0, 1, 2, 0},
				{5, 3, 4, 5, 3},
				{2, 0, 1, 2, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix, err := NewTorusMatrix(tt.width, tt.height)
			if err != nil {
				t.Fatalf("Failed to create matrix: %v", err)
			}

			hasher := NewMatrixHasher(matrix)
			wrapped := hasher.GenerateWrappedMatrix()

			if len(wrapped) != len(tt.expected) {
				t.Errorf("Expected %d rows, got %d", len(tt.expected), len(wrapped))
				return
			}

			for i, row := range wrapped {
				if len(row) != len(tt.expected[i]) {
					t.Errorf("Row %d: expected %d columns, got %d", i, len(tt.expected[i]), len(row))
					continue
				}

				for j, val := range row {
					if val != tt.expected[i][j] {
						t.Errorf("At position [%d][%d]: expected %d, got %d", i, j, tt.expected[i][j], val)
					}
				}
			}
		})
	}
}

func TestGenerateMatrixString(t *testing.T) {
	matrix, err := NewTorusMatrix(4, 4)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}

	hasher := NewMatrixHasher(matrix)
	matrixString := hasher.GenerateMatrixString()

	expectedLines := []string{
		"15,12,13,14,15,12",
		"3,0,1,2,3,0",
		"7,4,5,6,7,4",
		"11,8,9,10,11,8",
		"15,12,13,14,15,12",
		"3,0,1,2,3,0",
	}
	expected := strings.Join(expectedLines, "\n")

	if matrixString != expected {
		t.Errorf("Matrix string mismatch.\nExpected:\n%s\nGot:\n%s", expected, matrixString)
	}
}

func TestCalculateHash(t *testing.T) {
	matrix, err := NewTorusMatrix(4, 4)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}

	hasher := NewMatrixHasher(matrix)
	hash := hasher.CalculateHash()

	expected := "hJVz5fi5z2YecMNLsihGJQHBpAGUAYitNUOFGmjBg38="

	if hash != expected {
		t.Errorf("Hash mismatch.\nExpected: %s\nGot: %s", expected, hash)
	}
}

func TestValidateExpectedHash(t *testing.T) {
	matrix, err := NewTorusMatrix(4, 4)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}

	hasher := NewMatrixHasher(matrix)

	correctHash := "hJVz5fi5z2YecMNLsihGJQHBpAGUAYitNUOFGmjBg38="
	if err := hasher.ValidateExpectedHash(correctHash); err != nil {
		t.Errorf("Validation failed for correct hash: %v", err)
	}

	incorrectHash := "wrong_hash_value"
	if err := hasher.ValidateExpectedHash(incorrectHash); err == nil {
		t.Error("Expected validation to fail for incorrect hash")
	}
}

func TestHashConsistency(t *testing.T) {
	matrix, err := NewTorusMatrix(4, 4)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}

	hasher := NewMatrixHasher(matrix)

	hash1 := hasher.CalculateHash()
	hash2 := hasher.CalculateHash()
	hash3 := hasher.CalculateHash()

	if hash1 != hash2 || hash2 != hash3 {
		t.Errorf("Hash calculation is not consistent: %s, %s, %s", hash1, hash2, hash3)
	}
}
