package domain

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindNeighbors(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		index    int
		expected []int
	}{
		{
			name:     "4x4 matrix, index 5 (example from problem)",
			width:    4,
			height:   4,
			index:    5,
			expected: []int{0, 1, 2, 4, 6, 8, 9, 10},
		},
		{
			name:     "4x4 matrix, index 0 (top-left corner)",
			width:    4,
			height:   4,
			index:    0,
			expected: []int{15, 12, 13, 3, 1, 7, 4, 5},
		},
		{
			name:     "4x5 matrix, index 1",
			width:    5,
			height:   4,
			index:    1,
			expected: []int{15, 16, 17, 0, 2, 5, 6, 7},
		},
		{
			name:     "1x3 matrix, index 2 (wrapping test)",
			width:    3,
			height:   1,
			index:    2,
			expected: []int{1, 2, 0, 1, 0, 1, 2, 0},
		},
		{
			name:     "1x1 matrix, index 0 (all neighbors are self)",
			width:    1,
			height:   1,
			index:    0,
			expected: []int{0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix, err := NewTorusMatrix(tt.width, tt.height)
			if err != nil {
				t.Fatalf("Failed to create matrix: %v", err)
			}

			finder := NewNeighborFinder(matrix)
			neighbors, err := finder.FindNeighbors(tt.index)

			if err != nil {
				t.Fatalf("FindNeighbors failed: %v", err)
			}

			if !reflect.DeepEqual(neighbors, tt.expected) {
				t.Errorf("Expected neighbors %v, got %v", tt.expected, neighbors)
			}
		})
	}
}

func TestFindNeighborsInvalidIndex(t *testing.T) {
	matrix, _ := NewTorusMatrix(4, 4)
	finder := NewNeighborFinder(matrix)

	tests := []int{-1, 16, 100}

	for _, invalidIndex := range tests {
		_, err := finder.FindNeighbors(invalidIndex)
		if err == nil {
			t.Errorf("Expected error for invalid index %d", invalidIndex)
		}
	}
}

func TestAllDirections(t *testing.T) {
	expectedCount := 8
	if len(AllDirections) != expectedCount {
		t.Errorf("Expected %d directions, got %d", expectedCount, len(AllDirections))
	}

	directionSet := make(map[string]bool)
	for _, direction := range AllDirections {
		key := fmt.Sprintf("%d,%d", direction.RowOffset, direction.ColOffset)
		if directionSet[key] {
			t.Errorf("Duplicate direction found: %s", direction.Name)
		}
		directionSet[key] = true
	}

	expectedDirections := []NeighborDirection{
		{-1, -1, "TopLeft"},
		{-1, 0, "Top"},
		{-1, 1, "TopRight"},
		{0, -1, "Left"},
		{0, 1, "Right"},
		{1, -1, "BottomLeft"},
		{1, 0, "Bottom"},
		{1, 1, "BottomRight"},
	}

	if !reflect.DeepEqual(AllDirections, expectedDirections) {
		t.Errorf("AllDirections doesn't match expected order and values")
	}
}
