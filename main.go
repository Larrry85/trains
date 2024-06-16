package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Error: Incorrect number of command line arguments")
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

	movements, err := PlanTrainMovements(network, startStation, endStation, numTrains)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	for _, move := range movements {
		fmt.Println(strings.Join(move, " "))
	}
} // main() END

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
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



	return network, nil
} // ReadNetwork() END

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BFS to find all shortest paths from start to end
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
} // BFSAllPaths() END

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} // contains() END

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// PlanTrainMovements plans the movements for multiple trains
func PlanTrainMovements(network Network, start, end string, numTrains int) ([][]string, error) {
	allPaths, err := BFSAllPaths(network, start, end)
	if err != nil {
		return nil, err
	}

	paths := make([][]string, numTrains)
	for i := 0; i < numTrains; i++ {
		paths[i] = allPaths[i%len(allPaths)]
	}

	// Simulation of train movements
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
		movements = append(movements, move)
	}

	return movements, nil
} // PlanTrainMovements() END
