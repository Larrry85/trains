package A

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PriorityQueue []*Node

type Node struct {
	Station  string // The name of the station this node represents.
	Cost     int    // The cost to reach this node from the start node (typically the number of steps).
	Priority int    // The priority of this node in the priority queue, often based on cost plus a heuristic estimate.
	Parent   *Node  // A pointer to the parent node in the path, used to reconstruct the path after reaching the goal.
	TrainID  int    // The ID of the train that this node represents, ensuring that we track which train is at which node.
	Time     int    // The time step at which this node is occupied, used to avoid collisions by checking time and station occupancy.
}

type Station struct {
	name string
	x, y int
}

type Graph struct {
	stations    map[string]Station
	connections map[string][]string
}

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

	shortest := Paths(startStation, endStation, graph.connections)

	if len(shortest) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No path exists between the start and end stations")
		return
	}

	// Distribute trains across the shortest paths
	assignments := distributeTrainsInCycles(shortest, numTrains)
	// Simulate train movements across the shortest paths
	total := simulateTrainMovements(shortest, assignments)

	// Print the total movements
	fmt.Printf("Total Movements: %d\n", total)
	fmt.Println("***********")
}
func read(filePath string) (*Graph, error) {
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

// find routes using A* algorithm
func Paths(start, end string, graph map[string][]string) [][]string {

	numberOfPaths := 4
	trains := make([][]*Node, numberOfPaths)

	remainingPaths := numberOfPaths
	var result [][]string

	for turn := 0; remainingPaths > 0; turn++ {

		for i := 0; i < numberOfPaths; i++ {
			if len(trains[i]) > 0 && trains[i][len(trains[i])-1].Station == end {
				continue // Skip trains that have already reached the destination
			}

			path := cooperativeAStar(start, end, graph, i, trains, turn)

			if len(path) == 0 {
				break
			} else {
				trains[i] = path
				remainingPaths--
				var route []string
				// Collect movements for each turn starting from the first movement
				for _, node := range path {

					route = append(route, node.Station)
				}
				result = append(result, route)

			}
		}
	}

	return result
}

// A cooperative version of the A* algorithm that takes into consideration if stations are occupied/used or not.
func cooperativeAStar(start, end string, data map[string][]string, trainID int, trains [][]*Node, startTime int) []*Node {
	pq := &PriorityQueue{}
	heap.Init(pq)
	startNode := &Node{
		Station:  start,
		Cost:     0,
		Priority: heuristic(start, end),
		TrainID:  trainID,
		Time:     startTime,
	}
	heap.Push(pq, startNode)

	visited := make(map[string]bool)
	var finalPath []*Node

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Node)

		if current.Station == end {
			finalPath = reconstructPath(current)
			break
		}

		if visited[current.Station] {
			continue
		}
		visited[current.Station] = true

		for _, neighbor := range data[current.Station] {
			if isOccupied(neighbor, current.Time+1, trains) {
				continue
			}
			cost := current.Cost + 1
			priority := cost + heuristic(neighbor, end)
			neighborNode := &Node{
				Station:  neighbor,
				Cost:     cost,
				Priority: priority,
				Parent:   current,
				TrainID:  trainID,
				Time:     current.Time + 1,
			}
			heap.Push(pq, neighborNode)
		}
	}

	return finalPath
}

func isOccupied(station string, time int, trains [][]*Node) bool {
	for _, trainPath := range trains {

		for _, train := range trainPath {
			if train != nil && train.Station == station && train.Time == time {
				return true
			}
		}

	}
	return false
}

// ended up not using this because it's unnecessary. would be useful if the network was significantly bigger than the max limit of this program which is 10k stations
func heuristic(_, _ string) int {
	// Implement heuristic function (e.g., Manhattan distance, if applicable)
	return 1
}

func reconstructPath(node *Node) []*Node {
	var path []*Node
	for node != nil {
		path = append([]*Node{node}, path...)
		node = node.Parent
	}
	return path
}

// Implement heap.Interface for PriorityQueue
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*Node)
	*pq = append(*pq, node)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}
func distributeTrainsInCycles(paths [][]string, numTrains int) []int {
	numPaths := len(paths)
	trainAssignments := make([]int, numTrains)

	// Distribute trains across the paths in a balanced way
	for i := 0; i < numTrains; i++ {
		trainAssignments[i] = i % numPaths
	}

	/*/ Print the train assignments in cycles
	fmt.Println("Train assignments in cycles:")
	for i, assignment := range trainAssignments {
		fmt.Printf("Train %d assigned to Path %d: %v\n", i+1, assignment+1, paths[assignment])
	}*/
	fmt.Println()

	return trainAssignments
}

// Train represents a train with an ID and color.
type Train struct {
	ID    int
	Color string
}

// simulateTrainMovements simulates the movements of trains and prints their paths.
func simulateTrainMovements(paths [][]string, trainAssignments []int) int {
	numTrains := len(trainAssignments)
	trainPositions := make([]int, numTrains)
	stationQueues := make(map[string]*Queue)
	fmt.Print("Train movements:\n\n")

	// Initialize trains with colors
	trains := make([]Train, numTrains)
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
		trains[i] = Train{ID: i + 1, Color: color}
		trainPositions[i] = 0 // Initialize train positions
	}

	// Initialize station queues
	for i := 0; i < numTrains; i++ {
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

				// Check if this train is next in line at the current station
				if stationQueues[currentStation] != nil && stationQueues[currentStation].Front() == i {
					// Check if the next station is free from incoming trains
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
						moveLine = append(moveLine, fmt.Sprintf("\033[%smT%d\033[0m-%s", trains[i].Color, trains[i].ID, nextStation))

						// Update train's position and station queues
						trainPositions[i]++
						stationQueues[currentStation].Pop()

						// Remove train from the current station queue when it reaches the end station
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
	fmt.Print("\n")
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
