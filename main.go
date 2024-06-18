package main

import (
	"bufio"
<<<<<<< HEAD
=======
	"errors"
	"flag"
>>>>>>> origin/main
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

<<<<<<< HEAD
=======
// This Go program simulates the movement of multiple trains on a rail network, from a starting 
// station to an ending station, based on the shortest paths between the stations


// Station represents a train station with a name and coordinates
>>>>>>> origin/main
type Station struct {
	name string
	x, y int
}

type Graph map[string][]string

<<<<<<< HEAD
func networkMap(filename string) (Graph, map[string]Station, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
=======
// Network represents the entire rail network
// maps for stations and connections
type Network struct {
	Stations    map[string]Station
	Connections map[string][]string
}


// Train represents information about a train
type Train struct {
	ID       int
	Location string // Current location of the train
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {

	flag.Parse()

	if len(os.Args) != 5 { // too few or too many arguments
		fmt.Fprintln(os.Stderr, "\nError: Incorrect number of command line arguments")
		fmt.Fprint(os.Stderr, "Usage: go run main.go [path to file containing network map] [start station] [end station] [number of trains]\n\n")
		return
	}
	//  file path to the network data, the name of the starting station, the name of the ending station,
	// and the number of trains

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 { // no valid number of trains
		fmt.Fprintln(os.Stderr, "Error: Number of trains is not a valid positive integer")
		return
	}

	// calls the ReadNetwork function to parse the network data from the specified file
	network, err := ReadNetwork(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if _, exists := network.Stations[startStation]; !exists { // no start station
		fmt.Fprintln(os.Stderr, "Error: Start station does not exist")
		return
	}

	if _, exists := network.Stations[endStation]; !exists { // no end station
		fmt.Fprintln(os.Stderr, "Error: End station does not exist")
		return
	}

	if startStation == endStation { // same start and end
		fmt.Fprintln(os.Stderr, "Error: Start and end station are the same")
		return
	}
	
	// Check for stations with the same coordinates
    if err := checkDuplicateCoordinates(network); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

	movements, err := PlanTrainMovements(network, startStation, endStation, numTrains)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Print the train movements
	for _, move := range movements {
		fmt.Println(strings.Join(move, " "))

	}
	fmt.Print("\n***\n\n")
} // main() END


// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
}// checkDuplicateCoordinates() END


// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ReadNetwork reads and parses the network map from a file
func ReadNetwork(filePath string) (Network, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Network{}, err
>>>>>>> origin/main
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graph := make(Graph)
	stations := make(map[string]Station)
	var section string

	// Reads each line, identifying sections and parsing the station and connection data

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

<<<<<<< HEAD
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
=======
		switch section {
		case "stations": // in case of STATIONS
			parts := strings.Split(line, ",")
			if len(parts) != 3 { // must be a station and two coordinates
				return Network{}, errors.New("invalid station format")
			}
			name := strings.TrimSpace(parts[0]) // station name
			x, err1 := strconv.Atoi(strings.TrimSpace(parts[1])) // int coordinate
			y, err2 := strconv.Atoi(strings.TrimSpace(parts[2])) // int coordinate
			if err1 != nil || err2 != nil || x < 0 || y < 0 { // positive int coordinate
				return Network{}, errors.New("invalid station coordinates")
			}
			if _, exists := network.Stations[name]; exists { // cannot be two same stations 
				return Network{}, errors.New("duplicate station name")
			}
			// The line of code creates a new Station struct with the given name,
			//  x coordinate, and y coordinate, and then it stores this Station in 
			// the Stations map of the network under the key name. Essentially, it 
			// adds a new station to the rail network.
			network.Stations[name] = Station{name, x, y} // Add one station data to Station map
			// example: 
			// network.Stations["A"] = Station{"A", 10, 20}
			// example: 
			// map[string]Station{
			//    "A": Station{Name: "A", X: 10, Y: 20},
			//}


		case "connections": // in case of CONNECTIONS
			parts := strings.Split(line, "-")
			if len(parts) != 2 { // must be two stations and "-" between
				return Network{}, errors.New("invalid connection format")
			}
			from := strings.TrimSpace(parts[0]) // first station
			to := strings.TrimSpace(parts[1])   // second station
			if _, exists := network.Stations[from]; !exists { // if not found in STATIONS
				return Network{}, errors.New("connection references non-existent station")
			}
			if _, exists := network.Stations[to]; !exists { // if not found in STATIONS
				return Network{}, errors.New("connection references non-existent station")
			}
			if from == to { // if two same stations
				return Network{}, errors.New("connection references the same station")
			}
			for _, conn := range network.Connections[from] {
				if conn == to { //cannot be two same connections
					return Network{}, errors.New("duplicate connection")
				}
			}

			// reads connection data from a file and populates the 
			// Connections map in the Network struct. 

			// Adds first station, "from", to a Connections map
			network.Connections[from] = append(network.Connections[from], to)
			// Adds second station, "to", to a Connections map
			network.Connections[to] = append(network.Connections[to], from)
>>>>>>> origin/main
		}
		// example:
		// network.Connections["A"] = append(network.Connections["A"], "B")
		// network.Connections["B"] = append(network.Connections["B"], "A")
		// example:
		// map[string][]string{
		//    "A": {"B"},
		//    "B": {"A"},
		//}

	}

	if len(network.Stations) > 10000 {
		return Network{}, errors.New("too many stations")
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

<<<<<<< HEAD
	return graph, stations, nil
}

func Path(graph Graph, start, end string) [][]string {
=======

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Uses Breadth-First Search (BFS) to find all shortest paths
// from the start station to the end station
func BFSAllPaths(network Network, start, end string) ([][]string, error) {
>>>>>>> origin/main
	queue := [][]string{{start}}
	visited := make(map[string]bool)
	visited[start] = true
	var paths [][]string

	// Initializes a queue with the start station and iterates to find paths to the end station.
	// Stops when all shortest paths are found.
	// Ensures paths are unique and correct.

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

<<<<<<< HEAD
	networkFile := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintf(os.Stderr, "Error: invalid number of trains\n")
		os.Exit(1)
=======
	return paths, nil
} // BFSAllPaths() END


// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper function to check if a slice contains a specific string
// to avoid revisiting stations in the current path
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
>>>>>>> origin/main
	}

<<<<<<< HEAD
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
=======

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// chooseShortestPaths selects the shortest paths from a list of paths
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
} // chooseShortestPaths() END



///////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
Finding Shortest Paths: It finds all the shortest paths from the starting station to the ending station.
Assigning Paths: It assigns these paths to the trains evenly using a round-robin method.
Tracking Positions: It initializes the positions of all trains at the start station and simulates their movements step-by-step.
Generating Instructions: It generates and records movement instructions for each train at each step until all trains reach the end station.
*/

// Plans the movements of multiple trains from the start to the end station based on the shortest paths
func PlanTrainMovements(network Network, start, end string, numTrains int) ([][]string, error) {
	
	allPaths, err := BFSAllPaths(network, start, end) //  to get all shortest paths
	if err != nil {
		return nil, err
	} // returns a list of all such paths. If no paths are found, it returns an error


	// If there are multiple paths, choose the shortest ones
	allPaths = chooseShortestPaths(allPaths)

	//  assigns paths to the trains evenly using a round-robin method

	// initializes a slice paths to store the path each train will take
	paths := make([][]string, numTrains)
	// iterates over the number of trains and assigns paths to each train 
	// in a round-robin manner using the modulo operator (%)
	for i := 0; i < numTrains; i++ {
		paths[i] = allPaths[i%len(allPaths)]
	} // This ensures that if there are more trains than paths, 
	// paths will be reused and distributed evenly among the trains

	// Simulates Train Movements, Tracking Their Positions at Each Step
	// A map trainPositions is created to track the current position of each train.
	//  Initially, all trains are positioned at the start station
	movements := [][]string{}

	trainPositions := make(map[string]string)
	for i := 0; i < numTrains; i++ {
		trainPositions[fmt.Sprintf("T%d", i+1)] = start
	} // Each train is given a unique identifier (e.g., "T1", "T2", etc.)


	// calculates the maximum length of the paths assigned to trains (maxLen). 
	// This is the longest number of steps any train needs to reach the end station
	maxLen := 0
	for _, path := range paths {
		if len(path) > maxLen {
			maxLen = len(path)
		}
	}

	//  iterates from 1 to maxLen (starting from 1 because step 0 is the initial position)
	for i := 1; i < maxLen; i++ {
		// In each iteration, it creates a slice move to store the movements 
		// for the current step
		move := []string{}
		// iterates over each train and its assigned path. If the current step index i
		// is within the path length, it generates a movement instruction for the train 
		// o move to the next station in its path
		for t, path := range paths {
			if i < len(path) {
				train := fmt.Sprintf("T%d", t+1)
				nextStation := path[i]
				// The movement instructions are formatted as "trainID-nextStation" 
				// (e.g., "T1-B") and added to the move slice.
				move = append(move, fmt.Sprintf("%s-%s", train, nextStation))
				trainPositions[train] = nextStation
			}
		}
		if len(move) > 0 {
		// The move slice is then appended to the movements slice, 
		// which stores all movements step-by-step.
		movements = append(movements, move)
		}
>>>>>>> origin/main
	}

	fmt.Fprintf(os.Stderr, "Number of turns: %d\n", numTurns)
}
