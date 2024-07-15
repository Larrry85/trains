package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"stations/A"
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

	// Find all shortest paths between start and end station
	findPath := len(graph.connections)
	shortestPaths := findShortestPaths(findPath, graph, startStation, endStation)

	if len(shortestPaths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No path exists between the start and end stations")
		return
	}

	// Print all shortest paths found
	fmt.Println("Shortest path(s) found:")
	for i, path := range shortestPaths {
		fmt.Printf("Path %d: %v\n", i+1, path)
		fmt.Println()
	}

	// Distribute trains across the shortest paths
	trainAssignments := distributeTrainsInCycles(shortestPaths, numTrains)

	// Simulate train movements across the shortest paths
	totalMovements := simulateTrainMovements(shortestPaths, trainAssignments)

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

func findShortestPaths(findPath int, graph *Graph, start, end string) [][]string {
	if findPath >= 30 {
		return A.Paths(findPath, start, end, graph.connections)
	}
	// Dijkstra's algorithm to find all shortest paths from start to end
	type State struct {
		cost     int
		path     []string
		lastNode string
	}

	// Priority queue for Dijkstra's algorithm
	pq := priorityQueue{}
	heap.Init(&pq)

	// Map to store the shortest paths
	shortestPaths := [][]string{}
	visited := make(map[string]bool)

	// Push the initial state into the priority queue
	initialState := State{cost: 0, path: []string{start}, lastNode: start}
	heap.Push(&pq, &Item{value: &initialState, priority: 0})

	for pq.Len() > 0 {
		// Pop the state with the smallest cost
		current := heap.Pop(&pq).(*Item).value.(*State)

		// Check if we reached the end station
		if current.lastNode == end {
			shortestPaths = append(shortestPaths, current.path)
			continue
		}

		// Skip processing if the node has been visited
		if visited[current.lastNode] {
			continue
		}
		visited[current.lastNode] = true

		// Explore neighbors
		for _, neighbor := range graph.connections[current.lastNode] {
			if !visited[neighbor] {
				// Calculate the cost to reach the neighbor
				newCost := current.cost + 1 // Assuming each connection has equal cost (can be adjusted if needed)

				// Create a new state for the neighbor
				newPath := make([]string, len(current.path))
				copy(newPath, current.path)
				newPath = append(newPath, neighbor)

				newState := State{cost: newCost, path: newPath, lastNode: neighbor}

				// Push the new state into the priority queue
				heap.Push(&pq, &Item{value: &newState, priority: newCost})
			}
		}
	}

	return shortestPaths
}

type Item struct {
	value    interface{}
	priority int
}

type priorityQueue []*Item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].priority < pq[j].priority }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func distributeTrainsInCycles(paths [][]string, numTrains int) []int {
	numPaths := len(paths)
	trainAssignments := make([]int, numTrains)

	// Distribute trains across the paths in a balanced way
	for i := 0; i < numTrains; i++ {
		trainAssignments[i] = i % numPaths
	}

	// Print the train assignments in cycles
	fmt.Println("Train assignments in cycles:")
	for i, assignment := range trainAssignments {
		fmt.Printf("Train %d assigned to Path %d: %v\n", i+1, assignment+1, paths[assignment])
	}
	fmt.Println()

	return trainAssignments
}

func simulateTrainMovements(paths [][]string, trainAssignments []int) int {
	numTrains := len(trainAssignments)
	trainPositions := make([]int, numTrains)
	stationQueues := make(map[string]*Queue)

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
		if stationQueues[startStation] == nil {
			stationQueues[startStation] = NewQueue()
		}
		stationQueues[startStation].Push(i)
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
				if stationQueues[currentStation] != nil && stationQueues[currentStation].Front() == i {
					// Check if next station is free from incoming trains
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

						// Update train's position and station queues
						trainPositions[i]++
						stationQueues[currentStation].Pop()

						// Remove train from current station queue when it reaches end station
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

// Queue implementation for station queues
type Queue struct {
	items []int
}

func NewQueue() *Queue {
	return &Queue{items: []int{}}
}

func (q *Queue) Push(item int) {
	q.items = append(q.items, item)
}

func (q *Queue) Pop() {
	if len(q.items) == 0 {
		return
	}
	q.items = q.items[1:]
}

func (q *Queue) Front() int {
	if len(q.items) == 0 {
		return -1
	}
	return q.items[0]
}

func (q *Queue) Remove(val int) {
	for i := 0; i < len(q.items); i++ {
		if q.items[i] == val {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return
		}
	}
}
