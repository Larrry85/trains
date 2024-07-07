// path.go
package train

import (
	"container/heap"
	"fmt"
	"slices"
	"strings"
)

type Item struct {
	value    string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func buildAdjacencyList(connections Connections) map[string][]string {
	adjacencyList := make(map[string][]string)
	for _, connection := range connections {
		adjacencyList[connection.Start] = append(adjacencyList[connection.Start], connection.End)
		adjacencyList[connection.End] = append(adjacencyList[connection.End], connection.Start)
	}
	return adjacencyList
}

func FindShortestPath(start, end string, connections Connections) ([]string, bool) {
	adjacencyList := buildAdjacencyList(connections)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{value: start, priority: 0})

	distances := make(map[string]int)
	previous := make(map[string]string)

	for station := range adjacencyList {
		distances[station] = int(^uint(0) >> 1) // Max int
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

func ScheduleTrainMovements(start, end string, connections Connections, numTrains int) []string {
	var movements []string
	occupied := make(map[string]int)
	trains := []string{}
	trainPositions := make(map[string]string)

	for i := 0; i < numTrains; i++ {
		train := fmt.Sprintf("T%d", i+1)
		trains = append(trains, train)
		trainPositions[train] = start
	}
	spath, _ := FindShortestPath(start, end, connections)
	step := 0

	for !allTrainsReachedEnd(trainPositions, end) {
		trainsPaths := make(map[string][]string)
		var moves []string
		nextOccupied := make(map[string]int)

		for i, train := range trains {
			if trainPositions[train] != end {
				var path []string
				reachedDestinationOr1TurnAway := true
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
					allPaths, found := FindAllPaths(trainPositions[train], end, connections)
					if found {
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
		for station, count := range nextOccupied {
			occupied[station] = count
		}

		if len(moves) > 0 {
			movements = append(movements, strings.Join(moves, " "))
			step++
		}
	}
	return movements
}

func allTrainsReachedEnd(trainPositions map[string]string, end string) bool {
	for _, pos := range trainPositions {
		if pos != end {
			return false
		}
	}
	return true
}

func FindAllPaths(start, end string, connections Connections) ([][]string, bool) {
	adjacencyList := buildAdjacencyList(connections)
	var paths [][]string
	queue := [][]string{{start}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		last := path[len(path)-1]

		if last == end {
			paths = append(paths, path)
		}

		for _, neighbor := range adjacencyList[last] {
			if !contains(path, neighbor) {
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return paths, len(paths) > 0
}

func CountOverlap(path1, path2 []string) int {
	overlapCount := 0
	for i, s1 := range path1 {
		if i >= len(path2) {
			break
		}
		if s1 == path2[i] {
			overlapCount++
		}
	}
	return overlapCount
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
