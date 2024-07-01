package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	maxInt   = int(^uint(0) >> 1)
	maxTurns = 6
)

type Edge struct {
	to int
}

type Graph struct {
	nodes [][]Edge
}

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Error: incorrect number of command line arguments")
		return
	}

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrainsStr := os.Args[4]

	numTrains, err := strconv.Atoi(numTrainsStr)
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: number of trains is not a valid positive integer")
		return
	}

	stations, graph, err := parseNetworkMap(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	startIdx, ok := stations[startStation]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: start station does not exist")
		return
	}

	endIdx, ok := stations[endStation]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: end station does not exist")
		return
	}

	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: start and end stations are the same")
		return
	}

	dist, prev := dijkstra(graph, startIdx)

	if dist[endIdx] == maxInt {
		fmt.Fprintln(os.Stderr, "Error: no path between the start and end stations")
		return
	}

	path := reconstructPath(prev, startIdx, endIdx)
	if path == nil {
		fmt.Fprintln(os.Stderr, "Error: no path found")
		return
	}

	result := make([][]string, numTrains)
	for train := 0; train < numTrains; train++ {
		var moves []string
		var curr = startIdx
		turns := 0

		for i := 1; i < len(path); i++ {
			if turns >= maxTurns {
				break
			}

			moves = append(moves, fmt.Sprintf("T%d-%s", train+1, getKeyByValue(stations, path[i])))
			curr = path[i]
			turns++
		}

		result[train] = moves
	}

	for _, moves := range result {
		fmt.Println(strings.Join(moves, " "))
	}
}

func parseNetworkMap(filePath string) (map[string]int, *Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	stations := make(map[string]int)
	var graph Graph
	var idx int

	stationsSection := false
	connectionsSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "stations:") {
			stationsSection = true
			connectionsSection = false
			continue
		}
		if strings.HasPrefix(line, "connections:") {
			stationsSection = false
			connectionsSection = true
			continue
		}
		if strings.HasPrefix(line, "#") || line == "" {
			continue // Skip comments and empty lines
		}

		if stationsSection {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, nil, fmt.Errorf("Invalid station data format: %s", line)
			}
			name := strings.TrimSpace(parts[0])
			if _, found := stations[name]; found {
				return nil, nil, fmt.Errorf("Duplicate station name: %s", name)
			}
			stations[name] = idx
			idx++
			graph.nodes = append(graph.nodes, []Edge{})
		} else if connectionsSection {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("Invalid connection format: %s", line)
			}
			fromName := strings.TrimSpace(parts[0])
			toName := strings.TrimSpace(parts[1])
			fromIdx, ok := stations[fromName]
			if !ok {
				return nil, nil, fmt.Errorf("Connection refers to non-existent station: %s", fromName)
			}
			toIdx, ok := stations[toName]
			if !ok {
				return nil, nil, fmt.Errorf("Connection refers to non-existent station: %s", toName)
			}
			graph.nodes[fromIdx] = append(graph.nodes[fromIdx], Edge{toIdx})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return stations, &graph, nil
}

func getKeyByValue(m map[string]int, value int) string {
	for key, val := range m {
		if val == value {
			return key
		}
	}
	return ""
}

func dijkstra(graph *Graph, start int) ([]int, []int) {
	numNodes := len(graph.nodes)
	dist := make([]int, numNodes)
	prev := make([]int, numNodes)
	for i := range dist {
		dist[i] = maxInt
		prev[i] = -1
	}
	dist[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{value: start, priority: 0})

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item).value
		for _, edge := range graph.nodes[u] {
			v := edge.to
			alt := dist[u] + 1 // Assuming weight of 1 for each edge
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				heap.Push(&pq, &Item{value: v, priority: alt})
			}
		}
	}

	return dist, prev
}

func reconstructPath(prev []int, start, end int) []int {
	var path []int
	for u := end; u != -1; u = prev[u] {
		path = append([]int{u}, path...)
	}
	if path[0] != start {
		return nil
	}
	return path
}

type Item struct {
	value    int
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
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
