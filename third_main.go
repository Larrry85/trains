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

	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station cannot be the same")
		return
	}

	path, err := dijkstra(graph, startStation, endStation)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

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

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "stations:" {
			section = "stations"
		} else if line == "connections:" {
			section = "connections"
		} else if section == "stations" {
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

			for _, connected := range graph.connections[from] {
				if connected == to {
					return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
				}
			}

			graph.connections[from] = append(graph.connections[from], to)
			graph.connections[to] = append(graph.connections[to], from)
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

	stationOccupancy := make(map[string][]string) // Track trains at each station
	moveCount := 0

	for _, station := range path {
		moveLine := ""

		for i := 0; i < numTrains; i++ {
			// Move train to the next station if possible
			if canMoveTrain(stationOccupancy, station, trains[i]) {
				// Remove train from current station if it's moving from a different station
				for currentStation, trainsAtStation := range stationOccupancy {
					for j, train := range trainsAtStation {
						if train == trains[i] {
							stationOccupancy[currentStation] = append(trainsAtStation[:j], trainsAtStation[j+1:]...)
							break
						}
					}
				}

				// Add train to the new station
				stationOccupancy[station] = append(stationOccupancy[station], trains[i])

				// Prepare the line to print
				moveLine += fmt.Sprintf("%s-%s ", trains[i], station)
			}
		}

		if moveLine != "" {
			fmt.Println(strings.TrimSpace(moveLine))
			moveCount++
		}
	}

	fmt.Println() // Empty line after movements
	fmt.Printf("Movements: %d\n", moveCount)
	fmt.Println() // Empty line after movement count

	// Print "***********"
	fmt.Println("***********")
}

func canMoveTrain(stationOccupancy map[string][]string, station, train string) bool {
	// Check if the train is already at the destination station
	for _, t := range stationOccupancy[station] {
		if t == train {
			return false
		}
	}

	// Otherwise, train can move to the destination station
	return true
}