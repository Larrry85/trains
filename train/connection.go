package train

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ParseConnections reads connections from an input reader and parses them into a Connections slice.
func ParseConnections(r io.Reader) (Connections, error) {
	scanner := bufio.NewScanner(r)
	connections := Connections{}
	foundStationsSection := false
	foundConnectionsSection := false
	stations := make(map[string]bool)           // To track valid stations
	visitedConnections := make(map[string]bool) // To track visited connections
	duplicateErrors := make(map[string]bool)    // To track duplicate errors

	// Read each line from the input reader
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line starts with "stations:" or "connections:"
		if strings.HasPrefix(line, "stations:") {
			foundStationsSection = true
			continue // Skip the section header line
		} else if strings.HasPrefix(line, "connections:") {
			foundConnectionsSection = true
			continue // Skip the section header line
		}

		// If we have not found the "stations:" or "connections:" section yet, skip this line
		if !foundConnectionsSection {
			continue
		}

		parts := strings.Split(line, "-")

		// Ensure each line has exactly two parts (start and end stations)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		start, end := parts[0], parts[1]

		// Track valid stations in the "stations:" section
		if foundStationsSection {
			stations[start] = true
			stations[end] = true
		}

		// Check if stations referenced in connections exist in the "stations:" section
		if !stations[start] || !stations[end] {
			return nil, fmt.Errorf("station %s referenced in connections but not defined in stations section", start)
		}

		// Check for duplicate connection and reverse connection
		if visitedConnections[start+"-"+end] || visitedConnections[end+"-"+start] {
			// Only report the error if it hasn't been reported before
			if !duplicateErrors[start+"-"+end] {
				duplicateErrors[start+"-"+end] = true
				return nil, fmt.Errorf("duplicate route found between %s and %s", start, end)
			}
			continue // Skip adding to connections slice
		}

		visitedConnections[start+"-"+end] = true
		visitedConnections[end+"-"+start] = true

		// Add the connection to the Connections slice
		connections = append(connections, Connection{
			Start: start,
			End:   end,
		})
	}

	// Check if we found both the "stations:" and "connections:" sections
	if !foundStationsSection {
		return nil, errors.New("map file does not contain a 'stations:' section")
	}
	if !foundConnectionsSection {
		return nil, errors.New("map file does not contain a 'connections:' section")
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err // Return error if scanning fails
	}

	return connections, nil // Return parsed Connections slice
}
