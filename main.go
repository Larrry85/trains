package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"stations/train"
	"strconv"
	"strings"
)

type Station struct {
	name string
	x, y int
}

type Graph struct {
	stations    map[string]Station
	connections map[string][]string
}

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Error: Incorrect number of arguments")
		return
	}

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of trains must be a positive integer")
		return
	}

	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station cannot be the same")
		return
	}

	graph, err := readMap(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if _, exists := graph.stations[startStation]; !exists {
		fmt.Fprintln(os.Stderr, "Error: Start station does not exist")
		return
	}

	if _, exists := graph.stations[endStation]; !exists {
		fmt.Fprintln(os.Stderr, "Error: End station does not exist")
		return
	}

	shortestPath, _ := train.FindShortestPath(startStation, endStation, convertGraphToConnections(graph))

	if len(shortestPath) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No path exists between the start and end stations")
		return
	}

	// Schedule train movements
	trainMovements := train.ScheduleTrainMovements(startStation, endStation, convertGraphToConnections(graph), numTrains)

	// Print the movements
	fmt.Println("Train movements:")
	for _, move := range trainMovements {
		fmt.Println(move)
	}

	fmt.Printf("\nTotal number of movements: %d\n", len(trainMovements))
}

func readMap(filePath string) (*Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := &Graph{
		stations:    make(map[string]Station),
		connections: make(map[string][]string),
	}

	existingConnections := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	section := ""

	stationCount := 0
	connectionCount := 0
	stationsSectionExists := false
	connectionsSectionExists := false

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
			// Process station line
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

			if _, exists := graph.stations[name]; exists {
				return nil, fmt.Errorf("duplicate station name: %s", name)
			}

			for _, station := range graph.stations {
				if station.x == x && station.y == y {
					return nil, fmt.Errorf("duplicate coordinates for station %s", name)
				}
			}

			graph.stations[name] = Station{name, x, y}
			stationCount++
		} else if section == "connections" {
			// Process connection line
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid connection line: %s", line)
			}

			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])

			if _, exists := graph.stations[from]; !exists {
				return nil, fmt.Errorf("connection from non-existent station: %s", from)
			}

			if _, exists := graph.stations[to]; !exists {
				return nil, fmt.Errorf("connection to non-existent station: %s", to)
			}

			if from == to {
				return nil, fmt.Errorf("connection with same start and end station: %s", from)
			}

			connectionKey := fmt.Sprintf("%s-%s", from, to)
			reverseConnectionKey := fmt.Sprintf("%s-%s", to, from)
			if _, exists := existingConnections[connectionKey]; exists {
				return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
			}
			if _, exists := existingConnections[reverseConnectionKey]; exists {
				return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
			}

			existingConnections[connectionKey] = struct{}{}
			existingConnections[reverseConnectionKey] = struct{}{}

			graph.connections[from] = append(graph.connections[from], to)
			graph.connections[to] = append(graph.connections[to], from)
			connectionCount++
		}

		// Check station and connection count
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

	if len(graph.stations) == 0 {
		return nil, errors.New("map does not contain any stations")
	}

	if len(graph.connections) == 0 {
		return nil, errors.New("map does not contain any connections")
	}

	return graph, nil
}

func convertGraphToConnections(graph *Graph) train.Connections {
	connections := make(train.Connections, 0)

	for from, toList := range graph.connections {
		for _, to := range toList {
			connections = append(connections, train.Connection{
				Start: from,
				End:   to,
			})
		}
	}

	return connections
}
