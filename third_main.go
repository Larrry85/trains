package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
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

	// Check for same start and end station
	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station cannot be the same")
		return
	}

	// Find all paths between start and end station
	paths := findAllPaths(graph, startStation, endStation)

	// Print all routes found
	fmt.Println("All possible routes:")
	for i, path := range paths {
		fmt.Printf("Route %d: %v\n", i+1, path)
		fmt.Println()
	}

	// Distribute trains across all routes
	trainAssignments := distributeTrains(paths, numTrains)

	// Simulate train movements across all routes
	totalMovements := simulateTrainMovements(paths, trainAssignments)

	// Print the total movements
	fmt.Printf("Total Movements: %d\n", totalMovements)
	fmt.Println("***********")
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

	scanner := bufio.NewScanner(file)
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
			continue
		} else if line == "connections:" {
			section = "connections"
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

			for _, connected := range graph.connections[from] {
				if connected == to {
					return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
				}
			}

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

	if len(graph.stations) == 0 {
		return nil, errors.New("map does not contain any stations")
	}

	if len(graph.connections) == 0 {
		return nil, errors.New("map does not contain any connections")
	}

	return graph, nil
}

func findAllPaths(graph *Graph, start, end string) [][]string {
	visited := make(map[string]bool)
	var paths [][]string
	currentPath := []string{start}
	findPaths(graph, start, end, visited, currentPath, &paths)
	return paths
}

func findPaths(graph *Graph, current, end string, visited map[string]bool, currentPath []string, paths *[][]string) {
	visited[current] = true

	if current == end {
		*paths = append(*paths, append([]string{}, currentPath...))
	} else {
		for _, neighbor := range graph.connections[current] {
			if !visited[neighbor] {
				currentPath = append(currentPath, neighbor)
				findPaths(graph, neighbor, end, visited, currentPath, paths)
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
	}

	visited[current] = false
}

func distributeTrains(paths [][]string, numTrains int) []int {
	trainAssignments := make([]int, numTrains)

	// Calculate the number of trains to assign per route
	trainsPerRoute := numTrains / len(paths)
	extraTrains := numTrains % len(paths)

	// Prioritize paths by their length (shorter paths first)
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Assign trains to paths from shortest to longest
	for i := 0; i < len(paths); i++ {
		numTrainsForPath := trainsPerRoute
		if i < extraTrains {
			numTrainsForPath++
		}
		for j := 0; j < numTrainsForPath; j++ {
			trainAssignments[i*trainsPerRoute+j] = i
		}
	}

	return trainAssignments
}

func simulateTrainMovements(paths [][]string, trainAssignments []int) int {
	numTrains := len(trainAssignments)
	trainPositions := make([]int, numTrains)
	stationQueues := make(map[string][]int)

	// Initialize trains with colors
	trains := make([]string, numTrains)
	for i := 0; i < numTrains; i++ {
		switch i % 4 {
		case 0:
			trains[i] = fmt.Sprintf("\033[31mT%d\033[0m", i+1) // Red
		case 1:
			trains[i] = fmt.Sprintf("\033[33mT%d\033[0m", i+1) // Yellow
		case 2:
			trains[i] = fmt.Sprintf("\033[34mT%d\033[0m", i+1) // Blue
		case 3:
			trains[i] = fmt.Sprintf("\033[32mT%d\033[0m", i+1) // Green
		}
	}

	// Initialize train positions and station queues
	for i := range trainPositions {
		trainPositions[i] = 0
		startStation := paths[trainAssignments[i]][0]
		stationQueues[startStation] = append(stationQueues[startStation], i)
	}

	// Main simulation loop
	var steps int // Initialize steps counter
	for steps = 0; ; steps++ {
		var moveLine []string
		allTrainsAtEnd := true

		// Process each train
		for i := 0; i < numTrains; i++ {
			path := paths[trainAssignments[i]]
			if trainPositions[i] < len(path)-1 {
				allTrainsAtEnd = false
				currentStation := path[trainPositions[i]]
				nextStation := path[trainPositions[i]+1]

				// Check if this train is next in line at current station
				if len(stationQueues[currentStation]) > 0 && stationQueues[currentStation][0] == i {
					// Check if next station is free from incoming trains
					nextStationFree := true
					for _, trainIdx := range stationQueues[nextStation] {
						if trainPositions[trainIdx] == trainPositions[i]+1 {
							nextStationFree = false
							break
						}
					}

					if nextStationFree {
						moveLine = append(moveLine, fmt.Sprintf("%s-%s", trains[i], nextStation))

						// Update train's position and station queues
						trainPositions[i]++
						stationQueues[currentStation] = stationQueues[currentStation][1:]

						// Remove train from current station queue when it reaches end station
						if nextStation == path[len(path)-1] {
							stationQueues[nextStation] = stationQueues[nextStation][:0]
						} else {
							stationQueues[nextStation] = append(stationQueues[nextStation], i)
						}
					}
				}
			}
		}

		// Print movements if there are any
		if len(moveLine) > 0 {
			fmt.Println(strings.Join(moveLine, " "))
		}

		// Break loop if all trains have reached their destination
		if allTrainsAtEnd {
			break
		}

		// Add a check to prevent potential infinite loops
		if steps > 2*numTrains*len(paths[0]) {
			fmt.Fprintln(os.Stderr, "Error: Simulation exceeded maximum steps, possible infinite loop detected")
			return steps
		}
	}

	// Print movements count
	return steps
}
