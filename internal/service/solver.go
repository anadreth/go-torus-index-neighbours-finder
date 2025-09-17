package service

import (
	"fmt"
	"strconv"
	"strings"
	"torus-neighbors/internal/api"
	"torus-neighbors/internal/domain"

	"github.com/google/uuid"
)

type ChallengeResult struct {
	NeighborsString string
	MatrixHash      string
}

type TorusChallengeSolver struct {
	apiClient *api.Client
}

func NewTorusChallengeSolver(apiClient *api.Client) *TorusChallengeSolver {
	return &TorusChallengeSolver{
		apiClient: apiClient,
	}
}

func (s *TorusChallengeSolver) SolveChallenge(userIdentifier string) error {
	challengeUUID := uuid.New().String()
	fmt.Printf("Generated UUID for challenge: %s\n", challengeUUID)

	fmt.Println("Testing API connection...")
	if err := s.apiClient.Ping(); err != nil {
		return fmt.Errorf("failed to ping API: %w", err)
	}
	fmt.Println("API connection successful!")

	fmt.Println("Requesting challenge from API...")
	challenge, err := s.apiClient.GetChallenge(challengeUUID, userIdentifier)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}

	fmt.Printf("Received challenge: width=%s, height=%s, target_index=%s\n",
		challenge.SetX, challenge.SetY, challenge.SetZ)

	width, err := strconv.Atoi(challenge.SetX)
	if err != nil {
		return fmt.Errorf("invalid width value '%s': %w", challenge.SetX, err)
	}

	height, err := strconv.Atoi(challenge.SetY)
	if err != nil {
		return fmt.Errorf("invalid height value '%s': %w", challenge.SetY, err)
	}

	targetIndex, err := strconv.Atoi(challenge.SetZ)
	if err != nil {
		return fmt.Errorf("invalid target index value '%s': %w", challenge.SetZ, err)
	}

	fmt.Printf("Computing solution for %dx%d matrix, target index %d...\n", width, height, targetIndex)
	result, err := s.ComputeSolution(width, height, targetIndex)
	if err != nil {
		return fmt.Errorf("failed to compute solution: %w", err)
	}

	fmt.Printf("Solution computed:\n")
	fmt.Printf("  Neighbors: %s\n", result.NeighborsString)
	fmt.Printf("  Matrix Hash: %s\n", result.MatrixHash)

	fmt.Println("Submitting solution to API...", challenge.UUID, result.NeighborsString, result.MatrixHash)
	if err := s.apiClient.SubmitSolution(challenge.UUID, result.NeighborsString, result.MatrixHash); err != nil {
		return fmt.Errorf("failed to submit solution: %w", err)
	}

	fmt.Println("Challenge completed successfully!")
	return nil
}

func (s *TorusChallengeSolver) ComputeSolution(width, height, targetIndex int) (*ChallengeResult, error) {
	matrix, err := domain.NewTorusMatrix(width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to create torus matrix: %w", err)
	}

	if !matrix.IsValidIndex(targetIndex) {
		return nil, fmt.Errorf("target index %d is invalid for %dx%d matrix", targetIndex, width, height)
	}

	neighborFinder := domain.NewNeighborFinder(matrix)
	neighbors, err := neighborFinder.FindNeighbors(targetIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to find neighbors: %w", err)
	}

	neighborsStrings := make([]string, len(neighbors))
	for i, neighbor := range neighbors {
		neighborsStrings[i] = strconv.Itoa(neighbor)
	}
	neighborsString := strings.Join(neighborsStrings, ",")

	hasher := domain.NewMatrixHasher(matrix)
	matrixHash := hasher.CalculateHash()

	return &ChallengeResult{
		NeighborsString: neighborsString,
		MatrixHash:      matrixHash,
	}, nil
}

func (s *TorusChallengeSolver) ValidateLocalExample() error {
	fmt.Println("Validating implementation against known examples...")

	testCases := []struct {
		width       int
		height      int
		targetIndex int
		expected    []int
		description string
	}{
		{4, 4, 5, []int{0, 1, 2, 4, 6, 8, 9, 10}, "4x4 matrix, index 5"},
		{4, 4, 0, []int{15, 12, 13, 3, 1, 7, 4, 5}, "4x4 matrix, index 0"},
		{5, 4, 1, []int{15, 16, 17, 0, 2, 5, 6, 7}, "4x5 matrix, index 1"},
		{3, 1, 2, []int{1, 2, 0, 1, 0, 1, 2, 0}, "1x3 matrix, index 2"},
		{1, 1, 0, []int{0, 0, 0, 0, 0, 0, 0, 0}, "1x1 matrix, index 0"},
	}

	for _, tc := range testCases {
		result, err := s.ComputeSolution(tc.width, tc.height, tc.targetIndex)
		if err != nil {
			return fmt.Errorf("failed to compute solution for %s: %w", tc.description, err)
		}

		neighborStrings := strings.Split(result.NeighborsString, ",")
		computedNeighbors := make([]int, len(neighborStrings))
		for i, str := range neighborStrings {
			neighbor, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("failed to parse neighbor '%s': %w", str, err)
			}
			computedNeighbors[i] = neighbor
		}

		if len(computedNeighbors) != len(tc.expected) {
			return fmt.Errorf("%s: expected %d neighbors, got %d", tc.description, len(tc.expected), len(computedNeighbors))
		}

		for i, expected := range tc.expected {
			if computedNeighbors[i] != expected {
				return fmt.Errorf("%s: neighbor %d mismatch - expected %d, got %d",
					tc.description, i, expected, computedNeighbors[i])
			}
		}

		fmt.Printf("✓ %s: PASSED\n", tc.description)
	}

	fmt.Println("Validating hash calculation...")
	matrix, _ := domain.NewTorusMatrix(4, 4)
	hasher := domain.NewMatrixHasher(matrix)
	expectedHash := "hJVz5fi5z2YecMNLsihGJQHBpAGUAYitNUOFGmjBg38="

	if err := hasher.ValidateExpectedHash(expectedHash); err != nil {
		return fmt.Errorf("hash validation failed: %w", err)
	}
	fmt.Printf("✓ Hash validation: PASSED\n")

	fmt.Println("All local validations passed!")
	return nil
}
