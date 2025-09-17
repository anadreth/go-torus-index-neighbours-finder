package domain

import (
	"testing"
)

func TestNewTorusMatrix(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		height      int
		expectError bool
	}{
		{"Valid 4x4 matrix", 4, 4, false},
		{"Valid 1x1 matrix", 1, 1, false},
		{"Valid 3x5 matrix", 3, 5, false},
		{"Invalid zero width", 0, 4, true},
		{"Invalid zero height", 4, 0, true},
		{"Invalid negative width", -1, 4, true},
		{"Invalid negative height", 4, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix, err := NewTorusMatrix(tt.width, tt.height)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for width=%d, height=%d", tt.width, tt.height)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			w, h := matrix.Dimensions()
			if w != tt.width || h != tt.height {
				t.Errorf("Expected dimensions %dx%d, got %dx%d", tt.width, tt.height, w, h)
			}
		})
	}
}

func TestTorusMatrixTotalElements(t *testing.T) {
	tests := []struct {
		width    int
		height   int
		expected int
	}{
		{4, 4, 16},
		{1, 1, 1},
		{3, 5, 15},
		{10, 2, 20},
	}

	for _, tt := range tests {
		matrix, _ := NewTorusMatrix(tt.width, tt.height)
		actual := matrix.TotalElements()

		if actual != tt.expected {
			t.Errorf("For %dx%d matrix, expected %d elements, got %d",
				tt.width, tt.height, tt.expected, actual)
		}
	}
}

func TestIndexToCoordinates(t *testing.T) {
	matrix, _ := NewTorusMatrix(4, 4)

	tests := []struct {
		index       int
		expectedRow int
		expectedCol int
		expectError bool
	}{
		{0, 0, 0, false},  // Top-left
		{3, 0, 3, false},  // Top-right
		{15, 3, 3, false}, // Bottom-right
		{12, 3, 0, false}, // Bottom-left
		{5, 1, 1, false},  // Middle
		{16, 0, 0, true},  // Out of bounds
		{-1, 0, 0, true},  // Negative index
	}

	for _, tt := range tests {
		row, col, err := matrix.IndexToCoordinates(tt.index)

		if tt.expectError {
			if err == nil {
				t.Errorf("Expected error for index %d", tt.index)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for index %d: %v", tt.index, err)
			continue
		}

		if row != tt.expectedRow || col != tt.expectedCol {
			t.Errorf("For index %d, expected (%d,%d), got (%d,%d)",
				tt.index, tt.expectedRow, tt.expectedCol, row, col)
		}
	}
}

func TestCoordinatesToIndex(t *testing.T) {
	matrix, _ := NewTorusMatrix(4, 4)

	tests := []struct {
		row      int
		col      int
		expected int
	}{
		{0, 0, 0},  // Top-left
		{0, 3, 3},  // Top-right
		{3, 3, 15}, // Bottom-right
		{3, 0, 12}, // Bottom-left
		{1, 1, 5},  // Middle
		// Test wrapping behavior
		{-1, 0, 12},  // Wrap to bottom row
		{0, -1, 3},   // Wrap to rightmost column
		{4, 0, 0},    // Wrap to top row
		{0, 4, 0},    // Wrap to leftmost column
		{-1, -1, 15}, // Wrap both dimensions
		{5, 5, 5},    // Wrap both dimensions positive
	}

	for _, tt := range tests {
		actual := matrix.CoordinatesToIndex(tt.row, tt.col)

		if actual != tt.expected {
			t.Errorf("For coordinates (%d,%d), expected index %d, got %d",
				tt.row, tt.col, tt.expected, actual)
		}
	}
}

func TestIsValidIndex(t *testing.T) {
	matrix, _ := NewTorusMatrix(4, 4)

	tests := []struct {
		index    int
		expected bool
	}{
		{0, true},
		{15, true},
		{7, true},
		{16, false},
		{-1, false},
		{100, false},
	}

	for _, tt := range tests {
		actual := matrix.IsValidIndex(tt.index)

		if actual != tt.expected {
			t.Errorf("For index %d, expected %t, got %t", tt.index, tt.expected, actual)
		}
	}
}
