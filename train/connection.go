package train

import (
	"bufio"
	"io"
	"strings"
)

// Connection represents a direct connection between two stations.
type Connection struct {
	Start string // Starting station
	End   string // Ending station
}

// Connections is a slice of Connection, representing multiple train connections.
type Connections []Connection

// ParseConnections reads connections from an input reader and parses them into a Connections slice.
func ParseConnections(r io.Reader) (Connections, error) {
	scanner := bufio.NewScanner(r)
	connections := Connections{}

	// Read each line from the input reader
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")

		// Ensure each line has exactly two parts (start and end stations)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		// Add the connection to the Connections slice
		connections = append(connections, Connection{
			Start: parts[0],
			End:   parts[1],
		})
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err // Return error if scanning fails
	}

	return connections, nil // Return parsed Connections slice
}
