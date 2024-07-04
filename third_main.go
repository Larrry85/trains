package main

import (
	"bufio"
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

	// Run Dijkstra's algorithm
	path, err := dijkstra(graph, startStation, endStation)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Print train movements
	printTrainMovements(path, numTrains)
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

	scanner := bufio.NewScanner(file)
	section := ""

	stationCount := 0
	connectionCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "stations:" {
			section = "stations"
			continue
		} else if line == "connections:" {
			section = "connections"
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

			for _, connected := range graph.connections[from] {
				if connected == to {
					return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
				}
			}

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

	if len(graph.stations) == 0 {
		return nil, errors.New("map does not contain any stations")
	}

	if len(graph.connections) == 0 {
		return nil, errors.New("map does not contain any connections")
	}

	return graph, nil
}

func dijkstra(graph *Graph, start, end string) ([]string, error) {
	dist := make(map[string]int)
	prev := make(map[string]string)
	unvisited := make(map[string]bool)

	for station := range graph.stations {
		dist[station] = 1<<31 - 1 // equivalent to infinity
		unvisited[station] = true
	}
	dist[start] = 0

	for len(unvisited) > 0 {
		var u string
		minDist := 1<<31 - 1
		for station := range unvisited {
			if dist[station] < minDist {
				minDist = dist[station]
				u = station
			}
		}

		if u == end {
			break
		}

		delete(unvisited, u)

		for _, neighbor := range graph.connections[u] {
			if !unvisited[neighbor] {
				continue
			}

			alt := dist[u] + 1
			if alt < dist[neighbor] {
				dist[neighbor] = alt
				prev[neighbor] = u
			}
		}
	}

	path := []string{}
	for u := end; u != ""; u = prev[u] {
		path = append([]string{u}, path...)
	}

	if len(path) == 0 || path[0] != start {
		return nil, errors.New("no path exists between the start and end stations")
	}

	return path, nil
}

func printTrainMovements(path []string, numTrains int) {
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

	// Track trains waiting to enter each station
	stationQueues := make(map[string][]int)

	// Initialize train positions on the path
	trainPositions := make([]int, numTrains)
	for i := range trainPositions {
		trainPositions[i] = 0
		stationQueues[path[0]] = append(stationQueues[path[0]], i)
	}

	// Maximum steps to prevent infinite loops
	maxSteps := 1000
	var steps int

	// Main simulation loop
	for steps = 0; steps < maxSteps; steps++ {
		var moveLine []string
		allTrainsAtEnd := true

		// Process each train
		for i := 0; i < numTrains; i++ {
			if trainPositions[i] < len(path)-1 {
				allTrainsAtEnd = false
				currentStation := path[trainPositions[i]]
				nextStation := path[trainPositions[i]+1]

				// Check if this train is next in line at current station
				if len(stationQueues[currentStation]) > 0 && stationQueues[currentStation][0] == i {
					// Check if next station is free from incoming trains
					nextStationFree := true
					for _, trainIdx := range stationQueues[nextStation] {
						if trainPositions[trainIdx] == trainPositions[i]+1 {
							nextStationFree = false
							break
						}
					}

					if nextStationFree {
						moveLine = append(moveLine, fmt.Sprintf("%s-%s", trains[i], nextStation))

						// Update train's position and station queues
						trainPositions[i]++
						stationQueues[currentStation] = stationQueues[currentStation][1:]
						stationQueues[nextStation] = append(stationQueues[nextStation], i)
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
	}

	// Check if simulation exceeded maximum steps
	if steps >= maxSteps {
		fmt.Fprintln(os.Stderr, "Error: Simulation exceeded maximum number of steps, possible infinite loop detected")
		return
	}

	// Print movements count
	fmt.Println()
	fmt.Printf("Movements: %d\n", steps)
	fmt.Println()

	// Print "***********"
	fmt.Println("***********")
}
