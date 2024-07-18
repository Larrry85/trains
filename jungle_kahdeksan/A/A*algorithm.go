package A

import (
	"container/heap"
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

// find routes using A* algorithm
func Paths(filePath int, start, end string, graph map[string][]string) [][]string {

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