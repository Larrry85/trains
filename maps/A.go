package main

import (
	"bufio"
	"container/heap"
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

	// Find all shortest paths between start and end station
	shortestPaths := findShortestPaths(startStation, endStation, graph.connections)

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

type PriorityQueue []*Node

type Node struct {
	Station  string // The name of the station this node represents.
	Cost     int    // The cost to reach this node from the start node (typically the number of steps).
	Priority int    // The priority of this node in the priority queue, often based on cost plus a heuristic estimate.
	Parent   *Node  // A pointer to the parent node in the path, used to reconstruct the path after reaching the goal.
	TrainID  int    // The ID of the train that this node represents, ensuring that we track which train is at which node.
	Time     int    // The time step at which this node is occupied, used to avoid collisions by checking time and station occupancy.
}

// find routes using A* algorithm
func findShortestPaths(start, end string, graph map[string][]string) [][]string {
	numberOfPaths := 3
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
			// Check if the neighbor station is occupied at the next time step
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

func simulateTrainMovements(paths [][]string, trainAssignments []int, trains [][]*Node) int {
	numTrains := len(trainAssignments)
	trainPositions := make([]int, numTrains)
	stationQueues := make(map[string]*Queue)
	trainAtStation := make(map[string]map[int]bool) // Maps stations to a set of occupied times

	// Initialize train queues and occupancy maps
	for i := range trainPositions {
		trainPositions[i] = 0
		startStation := paths[trainAssignments[i]][0]
		if stationQueues[startStation] == nil {
			stationQueues[startStation] = NewQueue()
		}
		stationQueues[startStation].Push(i)

		if trainAtStation[startStation] == nil {
			trainAtStation[startStation] = make(map[int]bool)
		}
		trainAtStation[startStation][0] = true
	}

	steps := 0
	for {
		var moveLine []string
		allTrainsAtEnd := true
		nextTrainAtStation := make(map[string]map[int]bool) // For next step checks

		// Initialize the map for the next step
		for station := range trainAtStation {
			nextTrainAtStation[station] = make(map[int]bool)
		}

		// Process each train
		for i := 0; i < numTrains; i++ {
			path := paths[trainAssignments[i]]
			if trainPositions[i] < len(path)-1 {
				allTrainsAtEnd = false
				currentStation := path[trainPositions[i]]
				nextStation := path[trainPositions[i]+1]
				nextTime := steps + 1

				// Check if this train is next in line at the current station
				if stationQueues[currentStation] != nil && stationQueues[currentStation].Front() == i {
					// Check if the next station is free from incoming trains
					if !trainAtStation[nextStation][nextTime] {
						moveLine = append(moveLine, fmt.Sprintf("%s-%s", trains[i], nextStation))

						// Update train's position and station queues
						trainPositions[i]++
						stationQueues[currentStation].Pop()

						if nextTrainAtStation[nextStation] == nil {
							nextTrainAtStation[nextStation] = make(map[int]bool)
						}
						nextTrainAtStation[nextStation][nextTime] = true
					}
				}
			}
		}

		// Update the trainAtStation map for the next step
		for station, times := range nextTrainAtStation {
			trainAtStation[station] = times
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

		steps++
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
