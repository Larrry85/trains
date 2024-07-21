package pathfinder

import (
	"container/heap"
	"fmt"

	"strings"
)

// Item represents an element in the priority queue with a value, priority, and index.
type Item struct {
	value    string
	priority int
	index    int
}

// PriorityQueue implements a priority queue using a slice of Items.
type PriorityQueue []*Item

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less returns true if the priority of item i is less than the priority of item j.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

// Swap swaps the elements with indexes i and j in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item x to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Connection represents a connection between two stations with a travel time.
type Connection struct {
	Start string
	End   string
	Time  int
}

// Connections is a slice of Connection.
type Connections []Connection

// buildAdjacencyList builds an adjacency list from connections with travel times.
func buildAdjacencyList(connections Connections) map[string]map[string]int {
	adjacencyList := make(map[string]map[string]int)
	for _, connection := range connections {
		if adjacencyList[connection.Start] == nil {
			adjacencyList[connection.Start] = make(map[string]int)
		}
		if adjacencyList[connection.End] == nil {
			adjacencyList[connection.End] = make(map[string]int)
		}
		adjacencyList[connection.Start][connection.End] = connection.Time
		adjacencyList[connection.End][connection.Start] = connection.Time
	}
	return adjacencyList
}

/*/ Convert adjacencyList to map[string][]string
func convertAdjacencyList(adjacencyList map[string]map[string]int) map[string][]string {
	convertedAdjacencyList := make(map[string][]string)
	for station, neighbors := range adjacencyList {
		convertedAdjacencyList[station] = make([]string, 0, len(neighbors))
		for neighbor := range neighbors {
			convertedAdjacencyList[station] = append(convertedAdjacencyList[station], neighbor)
		}
	}
	return convertedAdjacencyList
}*/

// FindShortestPath finds the fastest path from start to end using Dijkstra's algorithm.
func FindShortestPath(start, end string, connections Connections) ([]string, bool) {
	// Initialize variables for Dijkstra's algorithm
	adjacencyList := buildAdjacencyList(connections)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{value: start, priority: 0})

	distances := make(map[string]int)
	previous := make(map[string]string)
	movements := make(map[string]int) // Track number of movements to each station

	// Initialize distances and movements to infinity
	for station := range adjacencyList {
		distances[station] = int(^uint(0) >> 1) // Infinity
		movements[station] = int(^uint(0) >> 1) // Infinity
	}
	distances[start] = 0
	movements[start] = 0

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		currentItem := heap.Pop(&pq).(*Item)
		currentStation := currentItem.value

		if currentStation == end {
			// Reconstruct the shortest path
			path := []string{}
			for at := end; at != ""; at = previous[at] {
				path = append([]string{at}, path...)
			}
			return path, true
		}

		if visited[currentStation] {
			continue
		}
		visited[currentStation] = true

		for neighbor, travelTime := range adjacencyList[currentStation] {
			if visited[neighbor] {
				continue
			}
			newDist := distances[currentStation] + travelTime
			newMovements := movements[currentStation] + 1 // Increment movements count

			if newMovements < movements[neighbor] {
				distances[neighbor] = newDist
				movements[neighbor] = newMovements
				previous[neighbor] = currentStation
				heap.Push(&pq, &Item{value: neighbor, priority: newDist})
			}
		}
	}

	// If no path found
	return nil, false
}

// ScheduleTrainMovements schedules the movements of multiple trains from start to end.
func ScheduleTrainMovements(start, end string, connections Connections, numTrains int) []string {
	// Convert Connections to a graph representation
	//graph := buildGraph(connections)

	var movements []string
	occupied := make(map[string]int)
	trains := make([]string, numTrains)
	trainPositions := make(map[string]string)

	// Initialize trains and their positions with color coding
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
		trainPositions[trains[i]] = start
	}

	// Precompute the fastest path using Dijkstra's algorithm
	fpath, _ := FindShortestPath(start, end, connections)

	step := 0
	maxSteps := 10000 // Limit steps to avoid infinite loop

	for !allTrainsReachedEnd(trainPositions, end) && step < maxSteps {

		trainsPaths := make(map[string][]string)
		var moves []string
		nextOccupied := make(map[string]int)

		for i, train := range trains {
			if trainPositions[train] != end {
				var path []string
				reachedDestinationOr1TurnAway := true

				// Check paths of other trains to determine optimal route for this train
				for j := 1; j <= i; j++ {
					if trainPositions[fmt.Sprintf("T%d", j)] != end {
						tPath := trainsPaths[fmt.Sprintf("T%d", j)]
						if len(tPath) > 1 && tPath[1] != end {
							reachedDestinationOr1TurnAway = false
						}
					}
				}

				if step == 0 && i == 0 {
					path = fpath
				} else if i == len(trains)-1 && reachedDestinationOr1TurnAway && len(fpath) == 2 {
					path = fpath
				} else {
					// Find all possible paths from current position to end
					allPaths, found := FindAllPaths(trainPositions[train], end, connections)
					if found {
						// Choose the best path based on overlap and other criteria
						for _, p := range allPaths {
							isGood := true
							for k := 1; k <= i; k++ {
								overlapCount := CountOverlap(trainsPaths[fmt.Sprintf("T%d", k)], p)
								if overlapCount > 4 {
									isGood = false
								}
							}
							isDuplicate := false
							for j := 1; j <= i; j++ {
								if slicesEqual(p, trainsPaths[fmt.Sprintf("T%d", j)]) {
									isDuplicate = true
								}
							}
							if nextOccupied[p[1]] == 0 && !contains(p[1:], start) && !isDuplicate && isGood {
								if len(p) > len(fpath)+2 {
									if contains(fpath, trainPositions[train]) {
										ind := slicesIndex(fpath, trainPositions[train])
										path = fpath[ind:]
									}
								} else {
									path = p
								}
								break
							}
						}
					}
				}

				if len(path) > 0 {
					trainsPaths[train] = path
					nextPos := path[1]
					if nextPos == end && trainPositions[train] != start || (nextOccupied[nextPos] == 0 && occupied[nextPos] == 0) {
						moves = append(moves, fmt.Sprintf("%s-%s", train, nextPos))
						if trainPositions[train] != start {
							occupied[trainPositions[train]]--
						}
						if nextPos != end {
							nextOccupied[nextPos]++
						}
						if nextPos == end {
							nextOccupied[nextPos] = 0
						}
						trainPositions[train] = nextPos
					}
				}
			}
		}

		// Update occupied stations
		for station, count := range nextOccupied {
			occupied[station] = count
		}

		// Add movements to result if there are any
		if len(moves) > 0 {
			movements = append(movements, strings.Join(moves, " "))
			step++
		}
	}

	fmt.Println("Train movements:")
	for _, move := range movements {
		fmt.Println(move)
	}

	fmt.Printf("\nTotal number of movements: %d\n", len(movements))
	fmt.Print("********************\n\n")
	return
}

// FindAllPaths finds all possible paths from start to end.
func FindAllPaths(start, end string, connections Connections) ([][]string, bool) {
	adjacencyList := buildAdjacencyList(connections)
	var paths [][]string
	queue := [][]string{{start}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		current := path[len(path)-1]

		if current == end {
			paths = append(paths, path)
		} else {
			for neighbor := range adjacencyList[current] {
				if !contains(path, neighbor) {
					newPath := append([]string(nil), path...)
					newPath = append(newPath, neighbor)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return paths, len(paths) > 0
}

// CountOverlap counts the number of overlapping elements between two slices.
func CountOverlap(a, b []string) int {
	count := 0
	set := make(map[string]bool)
	for _, v := range a {
		set[v] = true
	}
	for _, v := range b {
		if set[v] {
			count++
		}
	}
	return count
}

// contains checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// slicesEqual checks if two slices are equal.
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// slicesIndex returns the index of the first occurrence of item in slice.
func slicesIndex(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// allTrainsReachedEnd checks if all trains have reached the end station.
func allTrainsReachedEnd(trainPositions map[string]string, end string) bool {
	for _, pos := range trainPositions {
		if pos != end {
			return false
		}
	}
	return true
}

/*
func buildGraph(connections Connections) map[string][]string {
	graph := make(map[string][]string)
	for _, conn := range connections {
		graph[conn.Start] = append(graph[conn.Start], conn.End)
		graph[conn.End] = append(graph[conn.End], conn.Start) // If undirected
	}
	return graph
}*/
