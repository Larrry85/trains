//parser.go
package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	network "stations/go/network/dijkstra"
	"strconv"
	"strings"
)

func ParseConnections(r io.Reader) (network.Connections, error) {
	scanner := bufio.NewScanner(r)
	connections := network.Connections{}
	stations := make(map[string]network.Station)
	existingConnections := make(map[string]struct{})
	connectionsForStations := make(map[string]bool)
	stationsSectionExists := false
	connectionsSectionExists := false
	section := ""

	stationCount := 0
	connectionCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "stations:" {
			section = "stations"
			stationsSectionExists = true
			continue
		} else if line == "connections:" {
			section = "connections"
			connectionsSectionExists = true
			continue
		}

		if section == "stations" {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid station line: %s", line)
			}

			name := strings.TrimSpace(parts[0])
			x, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil || x < 0 {
				return nil, fmt.Errorf("invalid x coordinate for station %s", name)
			}
			y, err := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err != nil || y < 0 {
				return nil, fmt.Errorf("invalid y coordinate for station %s", name)
			}

			if _, exists := stations[name]; exists {
				return nil, fmt.Errorf("duplicate station name: %s", name)
			}
			for _, station := range stations {
				if station.X == x && station.Y == y {
					return nil, fmt.Errorf("duplicate coordinates for station %s", name)
				}
			}

			stations[name] = network.Station{Name: name, X: x, Y: y}
			connectionsForStations[name] = false 
			stationCount++
		} else if section == "connections" {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid connection line: %s", line)
			}

			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])

			startStation, exists := stations[from]
			if !exists {
				return nil, fmt.Errorf("connection from non-existent station: %s", from)
			}
			endStation, exists := stations[to]
			if !exists {
				return nil, fmt.Errorf("connection to non-existent station: %s", to)
			}
			if from == to {
				return nil, fmt.Errorf("connection with same start and end station: %s", from)
			}

			connectionKey := from + "-" + to
			reverseConnectionKey := to + "-" + from
			if _, exists := existingConnections[connectionKey]; exists {
				return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
			}
			if _, exists := existingConnections[reverseConnectionKey]; exists {
				return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
			}

			existingConnections[connectionKey] = struct{}{}
			existingConnections[reverseConnectionKey] = struct{}{}

			connections = append(connections, network.Connection{
				Start: startStation,
				End:   endStation,
			})
			connectionsForStations[from] = true // Mark that this station has a connection
			connectionsForStations[to] = true   // Mark that this station has a connection
			connectionCount++
		}

		if stationCount > 10000 {
			return nil, errors.New("map contains more than 10000 stations")
		}

		if connectionCount > 10000 {
			return nil, errors.New("map contains more than 10000 connections")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !stationsSectionExists {
		return nil, errors.New("map does not contain a \"stations:\" section")
	}

	if !connectionsSectionExists {
		return nil, errors.New("map does not contain a \"connections:\" section")
	}

	if len(stations) == 0 {
		return nil, errors.New("map does not contain any stations")
	}

	if len(connections) == 0 {
		return nil, errors.New("map does not contain any connections")
	}

	// Check if every station has at least one connection
	for station, hasConnection := range connectionsForStations {
		if !hasConnection {
			return nil, fmt.Errorf("no connection from station: %s", station)
		}
	}

	return connections, nil
}

func ReadMap(filePath string) (network.Connections, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseConnections(file)
}
