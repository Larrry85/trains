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

// ANSI color escape sequences
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Blue    = "\033[34m"
	Bold    = "\033[1m"
	Underline = "\033[4m"
)

// Station represents a train station with a name and coordinates
type Station struct {
	Name string
	X, Y int
}

// Connection represents a track between two stations
type Connection struct {
	From, To string
}

// Network represents the entire rail network
type Network struct {
	Stations    map[string]Station
	Connections map[string][]string
}

// Train represents information about a train
type Train struct {
	ID       int
	Location string // Current location of the train
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

	PrintRailwayMap(network)

	fmt.Print("\n")
	for _, move := range movements {
		fmt.Println(strings.Join(move, " "))
	}
	fmt.Print("\n***\n\n")
}

func checkDuplicateCoordinates(network Network) error {
	coordMap := make(map[string]string)
	for name, station := range network.Stations {
		coord := fmt.Sprintf("%d,%d", station.X, station.Y)
		if val, ok := coordMap[coord]; ok {
			return fmt.Errorf("error: Stations %s and %s have the same coordinates", name, val)
		}
		coordMap[coord] = name
	}
	return nil
}

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
	section := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "stations:" {
			section = "stations"
			continue
		}
		if line == "connections:" {
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

	return network, nil
}

func BFSAllPaths(network Network, start, end string) ([][]string, error) {
	queue := [][]string{{start}}
	paths := [][]string{}
	shortestPathLength := -1

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		last := path[len(path)-1]

		if last == end {
			if shortestPathLength == -1 || len(path) == shortestPathLength {
				shortestPathLength = len(path)
				paths = append(paths, path)
			} else if len(path) > shortestPathLength {
				break
			}
			continue
		}

		for _, neighbor := range network.Connections[last] {
			if !contains(path, neighbor) {
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func chooseShortestPaths(paths [][]string) [][]string {
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
	return shortestPaths
}

func PlanTrainMovements(network Network, start, end string, numTrains int) ([][]string, error) {
	allPaths, err := BFSAllPaths(network, start, end)
	if err != nil {
		return nil, err
	}

	allPaths = chooseShortestPaths(allPaths)

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
			if i < len(path) {
				train := fmt.Sprintf("T%d", t+1)
				nextStation := path[i]
				move = append(move, fmt.Sprintf("%s-%s", train, nextStation))
				trainPositions[train] = nextStation
			}
		}
		if len(move) > 0 {
			movements = append(movements, move)
		}
	}

	return movements, nil
}

func PrintRailwayMap(network Network) {
	maxX, maxY := 0, 0
	for _, station := range network.Stations {
		if station.X > maxX {
			maxX = station.X
		}
		if station.Y > maxY {
			maxY = station.Y
		}
	}

	grid := make([][]string, maxY+1)
	for i := range grid {
		grid[i] = make([]string, maxX+1)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	for _, station := range network.Stations {
		grid[station.Y][station.X] = Red + "X" + Reset
	}

	for from, connections := range network.Connections {
		fromStation := network.Stations[from]
		for _, to := range connections {
			toStation := network.Stations[to]
			if fromStation.X == toStation.X { // vertical connection
				for y := min(fromStation.Y, toStation.Y) + 1; y < max(fromStation.Y, toStation.Y); y++ {
					grid[y][fromStation.X] = Blue + "|" + Reset
				}
			} else if fromStation.Y == toStation.Y { // horizontal connection
				for x := min(fromStation.X, toStation.X) + 1; x < max(fromStation.X, toStation.X); x++ {
					grid[fromStation.Y][x] = Blue + "-" + Reset
				}
			} else { // diagonal or L-shaped connection
				if fromStation.X < toStation.X && fromStation.Y < toStation.Y {
					for x, y := fromStation.X+1, fromStation.Y+1; x < toStation.X && y < toStation.Y; x, y = x+1, y+1 {
						grid[y][x] = Blue + "+" + Reset
					}
				} else if fromStation.X > toStation.X && fromStation.Y < toStation.Y {
					for x, y := fromStation.X-1, fromStation.Y+1; x > toStation.X && y < toStation.Y; x, y = x-1, y+1 {
						grid[y][x] = Blue + "+" + Reset
					}
				} else if fromStation.X < toStation.X && fromStation.Y > toStation.Y {
					for x, y := fromStation.X+1, fromStation.Y-1; x < toStation.X && y > toStation.Y; x, y = x+1, y-1 {
						grid[y][x] = Blue + "+" + Reset
					}
				} else if fromStation.X > toStation.X && fromStation.Y > toStation.Y {
					for x, y := fromStation.X-1, fromStation.Y-1; x > toStation.X && y > toStation.Y; x, y = x-1, y-1 {
						grid[y][x] = Blue + "+" + Reset
					}
				}
			}
		}
	}

	// Print coordinates marking
	fmt.Printf("  ")
	for x := 0; x <= maxX; x++ {
		fmt.Printf(" %d", x)
	}
	fmt.Println()

	for y := 0; y <= maxY; y++ {
		fmt.Printf("%d ", y)
		for x := 0; x <= maxX; x++ {
			fmt.Print(grid[y][x])
		}
		fmt.Println()
	}

	fmt.Println()
	for _, station := range network.Stations {
		fmt.Printf("%s%s%s (%d, %d)\n", Red, station.Name, Reset, station.X, station.Y)
	}
}



func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
