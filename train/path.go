package train

import (
	"container/heap"
	"fmt"
	"slices"
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

// Connection represents a connection between two stations.
type Connection struct {
	Start string
	End   string
}

// Connections is a slice of Connection.
type Connections []Connection

// buildAdjacencyList builds an adjacency list from connections.
func buildAdjacencyList(connections Connections) map[string][]string {
	adjacencyList := make(map[string][]string)
	for _, connection := range connections {
		adjacencyList[connection.Start] = append(adjacencyList[connection.Start], connection.End)
		adjacencyList[connection.End] = append(adjacencyList[connection.End], connection.Start)
	}
	return adjacencyList
}

// FindShortestPath finds the shortest path from start to end using Dijkstra's algorithm.
func FindShortestPath(start, end string, connections Connections) ([]string, bool) {
	adjacencyList := buildAdjacencyList(connections)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{value: start, priority: 0})

	distances := make(map[string]int)
	previous := make(map[string]string)

	// Initialize distances to infinity
	for station := range adjacencyList {
		distances[station] = int(^uint(0) >> 1) // Infinity
	}
	distances[start] = 0

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		currentItem := heap.Pop(&pq).(*Item)
		currentStation := currentItem.value

		if currentStation == end {
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

		for _, neighbor := range adjacencyList[currentStation] {
			if visited[neighbor] {
				continue
			}
			newDist := distances[currentStation] + 1
			if newDist < distances[neighbor] {
				distances[neighbor] = newDist
				previous[neighbor] = currentStation
				heap.Push(&pq, &Item{value: neighbor, priority: newDist})
			}
		}
	}

	return nil, false
}

// ScheduleTrainMovements schedules the movements of multiple trains from start to end.
func ScheduleTrainMovements(start, end string, connections Connections, numTrains int) []string {
	var movements []string
	occupied := make(map[string]int)
	trains := make([]string, numTrains)
	trainPositions := make(map[string]string)

	// Initialize trains and their positions
	for i := 0; i < numTrains; i++ {
		train := fmt.Sprintf("T%d", i+1)
		trains[i] = train
		trainPositions[train] = start
	}

	// Precompute the shortest path using Dijkstra's algorithm
	spath, _ := FindShortestPath(start, end, connections)

	step := 0
	maxSteps := 100 // Limit steps to avoid infinite loop

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
					path = spath
				} else if i == len(trains)-1 && reachedDestinationOr1TurnAway && len(spath) == 2 {
					path = spath
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
								if slices.Equal(p, trainsPaths[fmt.Sprintf("T%d", j)]) {
									isDuplicate = true
								}
							}
							if nextOccupied[p[1]] == 0 && !contains(p[1:], start) && !isDuplicate && isGood {
								if len(p) > len(spath)+2 {
									if contains(spath, trainPositions[train]) {
										ind := slices.Index(spath, trainPositions[train])
										path = spath[ind:]
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

	return movements
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
			for _, neighbor := range adjacencyList[current] {
				if !contains(path, neighbor) {
					newPath := append([]string{}, path...)
					newPath = append(newPath, neighbor)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return paths, len(paths) > 0
}

// CountOverlap counts the number of overlapping stations between two paths.
func CountOverlap(path1, path2 []string) int {
	count := 0
	for _, station := range path1 {
		if contains(path2, station) {
			count++
		}
	}
	return count
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

// contains checks if a slice contains a specific element.
func contains(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

// slices package and functions like Equal and Index are assumed to be defined elsewhere
