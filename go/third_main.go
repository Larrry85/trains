package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
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
	paths := findShortestPaths(graph, startStation, endStation)
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No path exists between the start and end stations")
		return
	}

	// Print all routes found
	fmt.Println("Fastest routes:")
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

func findShortestPaths(graph *Graph, start, end string) [][]string {
	paths := [][]string{}
	queue := list.New()
	queue.PushBack([]string{start})
	visited := make(map[string]bool)

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]string)
		current := path[len(path)-1]

		if current == end {
			paths = append(paths, path)
			continue
		}

		if !visited[current] {
			visited[current] = true
			for _, neighbor := range graph.connections[current] {
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, neighbor)
				queue.PushBack(newPath)
			}
		}
	}

	return paths
}

func distributeTrains(paths [][]string, numTrains int) []int {
	trainAssignments := make([]int, numTrains)
	numPaths := len(paths)

	// Assign trains across all paths in a round-robin manner
	for i := 0; i < numTrains; i++ {
		trainAssignments[i] = i % numPaths
	}

	return trainAssignments
}

func simulateTrainMovements(paths [][]string, trainAssignments []int) int {
	numTrains := len(trainAssignments)
	trainPositions := make([]int, numTrains)
	stationQueues := make(map[string]*Queue)
	trains := make([]string, numTrains)
	for i := 0; i < numTrains; i++ {
		switch i % 4 {
		case 0:
			trains[i] = fmt.Sprintf("\033[31mT%d\033[0m", i+1)
		case 1:
			trains[i] = fmt.Sprintf("\033[33mT%d\033[0m", i+1)
		case 2:
			trains[i] = fmt.Sprintf("\033[34mT%d\033[0m", i+1)
		case 3:
			trains[i] = fmt.Sprintf("\033[32mT%d\033[0m", i+1)
		}
	}

	for i := range trainPositions {
		trainPositions[i] = 0
		startStation := paths[trainAssignments[i]][0]
		if stationQueues[startStation] == nil {
			stationQueues[startStation] = NewQueue()
		}
		stationQueues[startStation].Push(i)
	}

	var steps int
	for steps = 0; ; steps++ {
		var moveLine []string
		allTrainsAtEnd := true

		for i := 0; i < numTrains; i++ {
			path := paths[trainAssignments[i]]
			if trainPositions[i] < len(path)-1 {
				allTrainsAtEnd = false
				currentStation := path[trainPositions[i]]
				nextStation := path[trainPositions[i]+1]

				if stationQueues[currentStation] != nil && stationQueues[currentStation].Front() == i {
					nextStationFree := true
					if stationQueues[nextStation] != nil {
						for _, trainIndex := range stationQueues[nextStation].items {
							if trainPositions[trainIndex] == trainPositions[i]+1 {
								nextStationFree = false
								break
							}
						}
					}

					if nextStationFree {
						moveLine = append(moveLine, fmt.Sprintf("%s-%s", trains[i], nextStation))
						trainPositions[i]++
						stationQueues[currentStation].Pop()

						if nextStation == path[len(path)-1] {
							delete(stationQueues, nextStation)
						} else {
							if stationQueues[nextStation] == nil {
								stationQueues[nextStation] = NewQueue()
							}
							stationQueues[nextStation].Push(i)
						}
					}
				}
			}
		}

		if len(moveLine) > 0 {
			fmt.Println(strings.Join(moveLine, " "))
		}

		if allTrainsAtEnd {
			break
		}

		if steps > 2*numTrains*len(paths[0]) {
			fmt.Fprintln(os.Stderr, "Error: Simulation exceeded maximum steps, possible infinite loop detected")
			return steps
		}
	}

	return steps
}