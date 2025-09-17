package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"torus-neighbors/internal/api"
	"torus-neighbors/internal/service"
)

const (
	defaultAPIURL = "https://zadanie.openmed.sk"
	defaultUser   = ""
)

func main() {
	var (
		apiURL    = flag.String("api", defaultAPIURL, "API base URL")
		user      = flag.String("user", defaultUser, "User identifier")
		validate  = flag.Bool("validate", false, "Run local validation only (no API calls)")
		showUsage = flag.Bool("help", false, "Show usage information")
	)

	flag.Parse()

	if *showUsage {
		printUsage()
		return
	}

	// Initialize services
	apiClient := api.NewClient(*apiURL)
	solver := service.NewTorusChallengeSolver(apiClient)

	// Run local validation if requested
	if *validate {
		fmt.Println("Running local validation...")
		if err := solver.ValidateLocalExample(); err != nil {
			log.Fatalf("Local validation failed: %v", err)
		}
		fmt.Println("Local validation completed successfully!")
		return
	}

	// Run the full challenge
	fmt.Printf("Starting Torus Neighbors Challenge Solver\n")
	fmt.Printf("API URL: %s\n", *apiURL)
	fmt.Printf("User: %s\n", *user)
	fmt.Println()

	// First run local validation to ensure our implementation is correct
	fmt.Println("Running local validation before API interaction...")
	if err := solver.ValidateLocalExample(); err != nil {
		log.Fatalf("Local validation failed: %v", err)
	}
	fmt.Println()

	// Now solve the actual challenge from API
	if err := solver.SolveChallenge(*user); err != nil {
		log.Fatalf("Challenge failed: %v", err)
	}
}

func printUsage() {
	fmt.Printf(`Torus Neighbors Challenge Solver

This application solves the "Find the Cell's Neighbours" challenge by:
1. Creating a torus (wrap-around) matrix of specified dimensions
2. Finding all 8 neighbors of a given cell index
3. Computing the SHA256 hash of the extended matrix representation
4. Interacting with the challenge API to get problems and submit solutions

Usage:
  %s [options]

Options:
  -api <url>     API base URL (default: %s)
  -user <name>   User identifier for API requests (default: empty)
  -validate      Run local validation only, no API calls
  -help          Show this help message

Examples:
  # Run local validation only
  %s -validate
  
  # Solve challenge with default API
  %s -user "your-name"
  
  # Use custom API URL
  %s -api "https://custom-api.com" -user "your-name"

The application will:
1. Generate a UUID v4 for the attempt
2. Test API connectivity with /ping
3. Request a challenge from /challenge-me-easy
4. Compute the neighbors and matrix hash
5. Submit the solution back to the API

For more information about the problem, see the challenge description.
`, os.Args[0], defaultAPIURL, os.Args[0], os.Args[0], os.Args[0])
}
