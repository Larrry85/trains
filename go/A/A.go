package A

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"sort"
	astar "stations/go/network/astar"
	"strconv"
	"strings"
)

// PrintResult processes input and runs the simulation
func PrintResult() {
	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of trains must be a positive integer")
		return
	}

	graph, err := read(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	maxPaths := 8
	// Get distinct paths
	paths := findDistinctPaths(startStation, endStation, graph, maxPaths)

	// Distribute trains across paths
	trainAssignments := distributeTrainsAcrossPaths(paths, numTrains)

	// Simulate train movements
	simulateTrainMovements(paths, trainAssignments, startStation, endStation)

	fmt.Println("*******************")
}

// Read the input file and create a graph
func read(filePath string) (*astar.Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := &astar.Graph{
		Stations:    make(map[string]astar.Station),
		Connections: make(map[string][]string),
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

			if _, exists := graph.Stations[name]; exists {
				return nil, fmt.Errorf("duplicate station name: %s", name)
			}

			for _, station := range graph.Stations {
				if station.X == x && station.Y == y {
					return nil, fmt.Errorf("duplicate coordinates for station %s", name)
				}
			}

			graph.Stations[name] = astar.Station{
				Name: name,
				X:    x,
				Y:    y,
			}
			stationCount++
		} else if section == "connections" {
			// Process connection line
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid connection line: %s", line)
			}

			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])

			if _, exists := graph.Stations[from]; !exists {
				return nil, fmt.Errorf("connection from non-existent station: %s", from)
			}

			if _, exists := graph.Stations[to]; !exists {
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

			graph.Connections[from] = append(graph.Connections[from], to)
			graph.Connections[to] = append(graph.Connections[to], from)
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

	if len(graph.Stations) == 0 {
		return nil, errors.New("map does not contain any stations")
	}

	if len(graph.Connections) == 0 {
		return nil, errors.New("map does not contain any connections")
	}

	return graph, nil
}

// Find distinct paths between start and end
func findDistinctPaths(start, end string, graph *astar.Graph, maxPaths int) [][]string {
	paths := [][]string{}
	usedStations := make(map[string]struct{})
	allPaths := [][]string{}

	for len(paths) < maxPaths {
		path := aStarPathfinding(start, end, graph.Connections, usedStations)

		if len(path) == 0 {
			break
		}

		// Add the path to the list of all paths
		allPaths = append(allPaths, path)

		// Mark stations in this path as used
		for _, node := range path {
			if node != start && node != end {
				usedStations[node] = struct{}{}
			}
		}
	}

	// Sort paths by length (shortest first)
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Take up to maxPaths paths
	paths = allPaths
	if len(paths) > maxPaths {
		paths = paths[:maxPaths]
	}

	return paths
}

// A* pathfinding algorithm to find the optimal path
func aStarPathfinding(start, end string, connections map[string][]string, usedStations map[string]struct{}) []string {
	// Initialize the priority queue
	pq := &astar.PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &astar.Node{Station: start, Cost: 0, Priority: 0})

	// Map to keep track of the costs
	costSoFar := make(map[string]int)
	costSoFar[start] = 0

	// Map to keep track of the parent nodes for path reconstruction
	parentMap := make(map[string]string)

	for pq.Len() > 0 {
		currentNode := heap.Pop(pq).(*astar.Node)
		current := currentNode.Station

		// If we reached the goal, reconstruct the path
		if current == end {
			path := []string{}
			for current != start {
				path = append(path, current)
				current = parentMap[current]
			}
			path = append(path, start)

			// Reverse the path to get the correct order
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path
		}

		// Explore neighbors
		for _, neighbor := range connections[current] {
			if _, used := usedStations[neighbor]; used {
				continue
			}

			newCost := costSoFar[current] + 1 // Assuming all edges have a uniform cost
			if oldCost, ok := costSoFar[neighbor]; !ok || newCost < oldCost {
				costSoFar[neighbor] = newCost
				priority := newCost
				heap.Push(pq, &astar.Node{Station: neighbor, Cost: newCost, Priority: priority})
				parentMap[neighbor] = current
			}
		}
	}

	return []string{}
}

// Distribute trains across paths with specified constraints
func distributeTrainsAcrossPaths(paths [][]string, numTrains int) map[int]int {
	trainAssignments := make(map[int]int)
	pathUsageCount := make(map[int]int)

	// Sort paths by length
	sortedPaths := make([][]string, len(paths))
	copy(sortedPaths, paths)
	sort.Slice(sortedPaths, func(i, j int) bool {
		return len(sortedPaths[i]) < len(sortedPaths[j])
	})

	// Identify the shortest, second shortest, and longest paths
	shortestPathIndex := 0
	secondShortestPathIndex := 1
	longestPathIndex := len(paths) - 1

	// Assign trains
	for i := 0; i < numTrains; i++ {
		switch {
		case i < 5:
			trainAssignments[i] = shortestPathIndex
		case i < 8:
			trainAssignments[i] = secondShortestPathIndex
		case i < 9:
			trainAssignments[i] = longestPathIndex
		}
		pathUsageCount[trainAssignments[i]]++
	}

	return trainAssignments
}

// Simulate train movements on given paths
func simulateTrainMovements(paths [][]string, trainAssignments map[int]int, startStation, endStation string) int {
	totalMovements := 0
	numTrains := len(trainAssignments)
	trains := make([]astar.Train, numTrains)
	positions := make([]int, numTrains)
	completed := make([]bool, numTrains)
	trainLog := make([][]string, 0)
	occupiedStations := make(map[string]int)

	// Initialize trains
	for i := 0; i < numTrains; i++ {
		color := ""
		switch i % 4 {
		case 0:
			color = "31" // Red
		case 1:
			color = "33" // Yellow
		case 2:
			color = "34" // Blue
		case 3:
			color = "32" // Green
		}
		trains[i] = astar.Train{ID: i + 1, Color: color}
		occupiedStations[startStation]++
	}

	// Simulate movements
	for {
		allArrived := true
		turnLog := []string{}
		for i := 0; i < numTrains; i++ {
			if completed[i] {
				continue
			}

			pathIndex := trainAssignments[i]
			path := paths[pathIndex]

			// Check if the next position is free
			nextPosition := positions[i] + 1
			if nextPosition < len(path) {
				nextStation := path[nextPosition]
				if occupiedStations[nextStation] == 0 || nextStation == endStation {
					// Move the train
					if positions[i] > 0 {
						// Decrement the count of the current station only if it's not the start station
						occupiedStations[path[positions[i]]]--
					}
					positions[i]++
					occupiedStations[nextStation]++
					// Append the train movement with color
					turnLog = append(turnLog, fmt.Sprintf("\033[%smT%d-%s\033[0m", trains[i].Color, trains[i].ID, nextStation))
					totalMovements++
					allArrived = false

					if nextStation == endStation {
						completed[i] = true
					}
				}
			}
		}
		if len(turnLog) > 0 {
			trainLog = append(trainLog, turnLog)
		}
		if allArrived {
			break
		}
	}

	// Print the train movements
	fmt.Println("Train movements:")
	for _, turn := range trainLog {
		fmt.Println(strings.Join(turn, " "))
	}

	// Print the number of turns
	fmt.Printf("\nTotal Movements: %d\n", len(trainLog))

	return totalMovements
}
