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

type Node struct {
	Station  string
	Cost     int
	Priority int
	Parent   *Node
	TrainID  int
	Time     int
}

type PriorityQueue []*Node

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

	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station cannot be the same")
		return
	}

	// Find all shortest paths between start and end station
	shortestPaths := findShortestPaths(graph, startStation, endStation)

	if len(shortestPaths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No path exists between the start and end stations")
		return
	}

	// Distribute trains across the shortest paths
	trainAssignments := distributeTrainsInCycles(shortestPaths, numTrains)

	// Simulate train movements across the shortest paths
	totalMovements := simulateTrainMovements(shortestPaths, trainAssignments)

	// Print the total movements
	fmt.Printf("Total Movements: %d\n", totalMovements)
	fmt.Print("*************\n\n")
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

func findShortestPaths(graph *Graph, start, end string) [][]string {

	result1 := findPathsAStar(graph, start, end)
	result2 := findPathsDijkstra(graph, start, end)
	result3 := findPathsLimitedMoves(graph, start, end, 8)
	// Find the shortest path length in result1
	shortest1 := len(result1[0])
	for _, path := range result1 {
		if len(path) < shortest1 {
			shortest1 = len(path)
		}
	}

	// Find the shortest path length in result2
	shortest2 := len(result2[0])
	for _, path := range result2 {
		if len(path) < shortest2 {
			shortest2 = len(path)
		}
	}

	shortest3 := len(result3[0])
	for _, path := range result3 {
		if len(path) < shortest3 {
			shortest3 = len(path)
		}
	}

	// Compare the shortest path lengths
	if shortest1 <= shortest2 && shortest1 <= shortest3 {
		return result3
	} else if shortest2 <= shortest1 && shortest2 <= shortest3 {
		return result3
	}
	fmt.Println("mfklsdmfgösdg")
	return result3
}

func findPathsAStar(graph *Graph, start, end string) [][]string {
	numberOfPaths := 3
	trains := make([][]*Node, numberOfPaths)

	remainingPaths := numberOfPaths
	var result [][]string

	for turn := 0; remainingPaths > 0; turn++ {
		for i := 0; i < numberOfPaths; i++ {
			if len(trains[i]) > 0 && trains[i][len(trains[i])-1].Station == end {
				continue // Skip trains that have already reached the destination
			}

			path := cooperativeAStar(start, end, graph.connections, i, trains, turn)

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

func cooperativeAStar(start, end string, connections map[string][]string, trainID int, trains [][]*Node, startTime int) []*Node {
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

		for _, neighbor := range connections[current.Station] {
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

func heuristic(_, _ string) int {
	return 1 // Placeholder heuristic function
}

func reconstructPath(node *Node) []*Node {
	var path []*Node
	for node != nil {
		path = append([]*Node{node}, path...)
		node = node.Parent
	}
	return path
}

func findPathsDijkstra(graph *Graph, start, end string) [][]string {
	type State struct {
		cost     int
		path     []string
		lastNode string
	}
	pq := &priorityQueue{}
	heap.Init(pq)
	startState := &State{cost: 0, path: []string{start}, lastNode: start}
	heap.Push(pq, &Item{value: startState, priority: 0})

	visited := make(map[string]bool)
	shortestPaths := [][]string{}
	const maxPaths = 3 // Limit the number of routes to 3

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item).value.(*State)

		if current.lastNode == end {
			shortestPaths = append(shortestPaths, current.path)
			if len(shortestPaths) >= maxPaths {
				break // Stop searching after finding the maximum number of routes
			}
			continue
		}

		if visited[current.lastNode] {
			continue
		}
		visited[current.lastNode] = true

		for _, neighbor := range graph.connections[current.lastNode] {
			if !visited[neighbor] {
				newCost := current.cost + 1
				newPath := append([]string(nil), current.path...)
				newPath = append(newPath, neighbor)
				newState := &State{cost: newCost, path: newPath, lastNode: neighbor}
				heap.Push(pq, &Item{value: newState, priority: newCost})
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

func findPathsLimitedMoves(graph *Graph, start, end string, maxMovements int) [][]string {
	var result [][]string

	// Initialize the queue with the starting point
	type Path struct {
		nodes     []*Node
		movements int
	}
	queue := []Path{{nodes: []*Node{{Station: start}}, movements: 0}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		lastNode := currentPath.nodes[len(currentPath.nodes)-1]
		if lastNode.Station == end && currentPath.movements <= maxMovements {
			route := make([]string, len(currentPath.nodes))
			for i, node := range currentPath.nodes {
				route[i] = node.Station
			}
			result = append(result, route)
			continue
		}

		if currentPath.movements >= maxMovements {
			continue
		}

		for _, neighbor := range graph.connections[lastNode.Station] {
			newPath := Path{
				nodes:     append([]*Node(nil), currentPath.nodes...), // Copy the slice
				movements: currentPath.movements + 1,
			}
			newPath.nodes = append(newPath.nodes, &Node{Station: neighbor})
			stateIdentifier := generateStateIdentifier(newPath.nodes)

			if !visited[stateIdentifier] {
				visited[stateIdentifier] = true
				queue = append(queue, newPath)
			}
		}
	}

	if len(result) == 0 {
		fmt.Println("Error: No path exists between the start and end stations.")
	}

	return result
}

func generateStateIdentifier(nodes []*Node) string {
	var identifier strings.Builder
	for _, node := range nodes {
		identifier.WriteString(node.Station)
		identifier.WriteString("-")
	}
	return identifier.String()
}
