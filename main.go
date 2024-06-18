package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Station struct {
	name string
	x, y int
}

type Graph map[string][]string

func networkMap(filename string) (Graph, map[string]Station, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graph := make(Graph)
	stations := make(map[string]Station)
	var section string

	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), "#")[0] // Remove comments
		if line == "" {
			continue // Skip empty lines
		}

		if line == "stations:" {
			section = "stations"
			continue
		} else if line == "connections:" {
			section = "connections"
			continue
		}

		if section == "stations" {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, nil, fmt.Errorf("Error: invalid station line format")
			}
			name := strings.TrimSpace(parts[0])
			x, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, nil, fmt.Errorf("Error: invalid X-coordinate for station %s", name)
			}
			y, err := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err != nil {
				return nil, nil, fmt.Errorf("Error: invalid Y-coordinate for station %s", name)
			}
			stations[name] = Station{name, x, y}
		} else if section == "connections" {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("Error: invalid connection line format")
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			if _, ok := graph[from]; !ok {
				graph[from] = make([]string, 0)
			}
			graph[from] = append(graph[from], to)
			if _, ok := graph[to]; !ok {
				graph[to] = make([]string, 0)
			}
			graph[to] = append(graph[to], from)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return graph, stations, nil
}

func Path(graph Graph, start, end string) [][]string {
	queue := [][]string{{start}}
	visited := make(map[string]bool)
	visited[start] = true
	var paths [][]string

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node == end {
			paths = append(paths, path)
		}

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return paths
}

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Error: incorrect number of arguments\n")
		os.Exit(1)
	}

	networkFile := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintf(os.Stderr, "Error: invalid number of trains\n")
		os.Exit(1)
	}

	graph, stations, err := networkMap(networkFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if _, ok := stations[startStation]; !ok {
		fmt.Fprintf(os.Stderr, "Error: start station %s does not exist\n", startStation)
		os.Exit(1)
	}

	if _, ok := stations[endStation]; !ok {
		fmt.Fprintf(os.Stderr, "Error: end station %s does not exist\n", endStation)
		os.Exit(1)
	}

	if startStation == endStation {
		fmt.Fprintf(os.Stderr, "Error: start and end stations cannot be the same\n")
		os.Exit(1)
	}

	paths := Path(graph, startStation, endStation)

	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no path found between %s and %s\n", startStation, endStation)
		os.Exit(1)
	}

	// Sort paths by length (number of stations)
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Output train movements
	numTurns := 0
	for turn := 0; turn < len(paths[0]); turn++ {
		for trainID := 1; trainID <= numTrains; trainID++ {
			for _, path := range paths {
				if trainID > len(path) {
					continue
				}
				if turn < len(path) {
					fmt.Printf("T%d-%s ", trainID, path[turn])
				}
			}
		}
		fmt.Println()
		numTurns++
	}

	fmt.Fprintf(os.Stderr, "Number of turns: %d\n", numTurns)
}
