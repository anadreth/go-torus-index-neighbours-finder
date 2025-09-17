# Torus Neighbors Challenge Solver

A Go application that solves the "Find the Cell's Neighbours" challenge using clean code practices and SOLID principles.

## Problem Description

Given positive integers `h` (height), `w` (width) and a non-negative integer `i` (target index), the program finds all indices surrounding `i` in a torus matrix.

A torus matrix is a wrap-around matrix where:
- Elements are numbered sequentially from 0 to (h×w-1)
- Top edge connects to bottom edge
- Left edge connects to right edge
- Each cell has exactly 8 neighbors (including diagonals)

### Example
For a 4×4 matrix (h=4, w=4):
```
Original:     Extended (showing wrapping):
0  1  2  3    15 12 13 14 15 12
4  5  6  7 -> 3  0  1  2  3  0
8  9  10 11   7  4  5  6  7  4
12 13 14 15   11 8  9  10 11 8
              15 12 13 14 15 12
              3  0  1  2  3  0
```

For index 5, neighbors are: `0,1,2,4,6,8,9,10`

## Architecture

The solution follows clean architecture principles with clear separation of concerns:

```
cmd/                    # Application entry point
├── main.go            # CLI and application bootstrap

internal/
├── domain/            # Core business logic (domain layer)
│   ├── matrix.go      # Torus matrix representation
│   ├── neighbors.go   # Neighbor finding algorithm
│   ├── hasher.go      # Matrix hashing functionality
│   └── *_test.go      # Unit tests
│
├── service/           # Application services (use case layer)
│   └── solver.go      # Challenge orchestration
│
└── api/               # External API integration (infrastructure layer)
    ├── client.go      # HTTP API client
    └── client_test.go # Integration tests
```

### Design Principles Applied

1. **Single Responsibility Principle**: Each class/struct has one reason to change
2. **Open/Closed Principle**: Code is open for extension, closed for modification
3. **Dependency Inversion**: High-level modules don't depend on low-level modules
4. **Interface Segregation**: Small, focused interfaces
5. **Clean Code**: Descriptive names, short functions, clear intent

### Key Components

- **`TorusMatrix`**: Core domain entity representing the matrix with coordinate transformations
- **`NeighborFinder`**: Service responsible for finding cell neighbors
- **`MatrixHasher`**: Handles SHA256 hash calculation of extended matrix
- **`TorusChallengeSolver`**: Orchestrates the complete solution workflow
- **`Client`**: HTTP client for API interactions

## Usage

### Prerequisites
- Go 1.21 or later
- Internet connection for API challenges

### Building
```bash
make build
```

### Running Local Validation
Test the implementation against known examples:
```bash
make validate
```

### Solving API Challenge
```bash
make solve
# or directly:
./bin/torus-neighbors -user "your-name"
```

### Running Tests
```bash
make test                # Run all tests
make test-coverage      # Generate coverage report
```

### Available Commands
```bash
make help               # Show all available commands
make build              # Build the application
make test               # Run unit tests
make validate           # Run local validation
make solve              # Solve challenge interactively
make lint               # Run code linters
make clean              # Clean build artifacts
```

## API Integration

The application interacts with the challenge API:

1. **Ping**: `GET /ping` - Test connectivity
2. **Get Challenge**: `POST /challenge-me-easy` - Request a new challenge
3. **Submit Solution**: `POST /challenge-me-easy` - Submit computed solution

### Request/Response Format

**Challenge Request:**
```json
{
  "uuid": "generated-uuid-v4",
  "user": "your-identifier"
}
```

**Challenge Response:**
```json
{
  "uuid": "generated-uuid-v4",
  "set_x": "4",    // width
  "set_y": "4",    // height
  "set_z": "5"     // target index
}
```

**Solution Submission:**
```json
{
  "uuid": "generated-uuid-v4",
  "result": "0,1,2,4,6,8,9,10",
  "hash": "hJVz5fi5z2YecMNLsihGJQHBpAGUAYitNUOFGmjBg38="
}
```

## Algorithm Details

### Neighbor Finding
1. Convert linear index to 2D coordinates (row, col)
2. For each of 8 directions, calculate neighbor coordinates
3. Apply modular arithmetic for torus wrapping
4. Convert back to linear indices

### Hash Calculation
1. Generate extended matrix with wrapped borders
2. Create comma-separated string representation
3. Compute SHA256 hash
4. Encode to base64

### Complexity
- Time: O(1) for neighbor finding, O(w×h) for hash calculation
- Space: O(w×h) for extended matrix representation

## Testing

Comprehensive test coverage includes:
- Unit tests for all domain logic
- Property-based tests for coordinate transformations  
- Integration tests for API client
- Validation against known examples

Run with coverage:
```bash
make test-coverage
open coverage.html
```

## Error Handling

Robust error handling throughout:
- Input validation with descriptive error messages
- Network error handling with retries
- Graceful degradation for API failures
- Clear error propagation up the call stack

## Contributing

1. Follow Go conventions and gofmt styling
2. Write tests for new functionality
3. Update documentation as needed
4. Run `make check` before submitting