package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Station represents a train station with a name and coordinates
type Station struct {
	name string
	X    int // X coordinate of the station
	Y    int // Y coordinate of the station
}

// Network represents the entire rail network
type Network struct {
	Stations    map[string]Station  // Map of station name to Station struct
	Connections map[string][]string // Map of station name to list of connected station names
}

func main() {
	flag.Parse()

	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "\nError: Incorrect number of command line arguments")
		fmt.Fprint(os.Stderr, "Usage: go run main.go [path to file containing network map] [start station] [end station] [number of trains]\n\n")
		return
	}

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of trains is not a valid positive integer")
		return
	}

	network, err := ReadNetwork(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if _, exists := network.Stations[startStation]; !exists {
		fmt.Fprintln(os.Stderr, "Error: Start station does not exist")
		return
	}

	if _, exists := network.Stations[endStation]; !exists {
		fmt.Fprintln(os.Stderr, "Error: End station does not exist")
		return
	}

	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station are the same")
		return
	}

	if err := checkDuplicateCoordinates(network); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	movements, err := PlanTrainMovements(network, startStation, endStation, numTrains)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	for _, move := range movements {
		fmt.Println(strings.Join(move, " "))
	}
	fmt.Print("\n***\n\n")
}

// Function to check for stations with the same coordinates
func checkDuplicateCoordinates(network Network) error {
	coordMap := make(map[string]string) // Map to store coordinates as "X,Y" -> "station1,station2"
	for name, station := range network.Stations {
		coord := fmt.Sprintf("%d,%d", station.X, station.Y)
		if val, ok := coordMap[coord]; ok {
			return fmt.Errorf("error: Stations %s and %s have the same coordinates", name, val)
		}
		coordMap[coord] = name
	}
	return nil
}

// ReadNetwork reads and parses the network map from a file
func ReadNetwork(filePath string) (Network, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Network{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	network := Network{
		Stations:    make(map[string]Station),
		Connections: make(map[string][]string),
	}
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

		switch section {
		case "stations":
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return Network{}, errors.New("invalid station format")
			}
			name := strings.TrimSpace(parts[0])
			x, err1 := strconv.Atoi(strings.TrimSpace(parts[1]))
			y, err2 := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err1 != nil || err2 != nil || x < 0 || y < 0 {
				return Network{}, errors.New("invalid station coordinates")
			}
			if _, exists := network.Stations[name]; exists {
				return Network{}, errors.New("duplicate station name")
			}
			network.Stations[name] = Station{name, x, y}

		case "connections":
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return Network{}, errors.New("invalid connection format")
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			if _, exists := network.Stations[from]; !exists {
				return Network{}, errors.New("connection references non-existent station")
			}
			if _, exists := network.Stations[to]; !exists {
				return Network{}, errors.New("connection references non-existent station")
			}
			if from == to {
				return Network{}, errors.New("connection references the same station")
			}
			for _, conn := range network.Connections[from] {
				if conn == to {
					return Network{}, errors.New("duplicate connection")
				}
			}
			network.Connections[from] = append(network.Connections[from], to)
			network.Connections[to] = append(network.Connections[to], from)
		}
	}

	if len(network.Stations) > 10000 {
		return Network{}, errors.New("too many stations")
	}

	if err := scanner.Err(); err != nil {
		return Network{}, err
	}

	return network, nil
}

// Uses Breadth-First Search (BFS) to find all shortest paths from the start station to the end station
func BFSAllPaths(network Network, start, end string) ([][]string, error) {
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

		for _, neighbor := range network.Connections[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	if len(paths) == 0 {
		return nil, errors.New("no path found")
	}

	return paths, nil
}

// Helper function to select the shortest paths from a list of paths
func chooseShortestPaths(paths [][]string, start, end string, numTrains int) [][]string {
	shortestLength := len(paths[0])
	shortestPaths := [][]string{paths[0]}
	for _, path := range paths[1:] {
		if len(path) < shortestLength {
			shortestLength = len(path)
			shortestPaths = [][]string{path}
		} else if len(path) == shortestLength {
			shortestPaths = append(shortestPaths, path)
		}
	}
	// Ensure the number of shortest paths matches the number of trains
	if len(shortestPaths) < numTrains {
		for i := len(shortestPaths); i < numTrains; i++ {
			shortestPaths = append(shortestPaths, shortestPaths[i%len(shortestPaths)])
		}
	}
	return shortestPaths
}

// Plans the movements of multiple trains from the start to the end station based on the shortest paths
func PlanTrainMovements(network Network, start, end string, numTrains int) ([][]string, error) {
	allPaths, err := BFSAllPaths(network, start, end)
	if err != nil {
		return nil, err
	}

	allPaths = chooseShortestPaths(allPaths, start, end, numTrains)

	paths := make([][]string, numTrains)
	for i := 0; i < numTrains; i++ {
		paths[i] = allPaths[i%len(allPaths)]
	}

	movements := [][]string{}
	trainPositions := make(map[string]string)
	for i := 0; i < numTrains; i++ {
		trainPositions[fmt.Sprintf("T%d", i+1)] = start
	}

	maxLen := 0
	for _, path := range paths {
		if len(path) > maxLen {
			maxLen = len(path)
		}
	}

	for i := 1; i < maxLen; i++ {
		move := []string{}
		for t, path := range paths {
			train := fmt.Sprintf("T%d", t+1)
			if i < len(path) {
				nextStation := path[i]
				move = append(move, fmt.Sprintf("%s-%s", train, nextStation))
				trainPositions[train] = nextStation
			} else {
				move = append(move, fmt.Sprintf("%s-%s", train, trainPositions[train]))
			}
		}
		movements = append(movements, move)
	}

	return movements, nil
}
