package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Station struct {
	Name     string
	X, Y     int            // Example: coordinates, not used in this simplified approach
	Adjacent map[string]int // Adjacent stations and distances (weights)
}

var stations map[string]*Station

func main() {
	if len(os.Args) < 5 {
		fmt.Println("usage example: go run main.go london.txt waterloo victoria 5")
		return
	}

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	trainCount := os.Args[4]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inStations := false
	startFound := false
	endFound := false

	var connections []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, "stations:") {
			inStations = true
			continue
		}

		if inStations {
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				stationName := parts[0]
				if stationName == startStation {
					fmt.Printf("start station: %s\n", parts[0])
					startFound = true
				} else if stationName == endStation {
					fmt.Printf("end station: %s\n", parts[0])
					endFound = true
				}
			}
		}

		if strings.HasPrefix(line, "connections") {
			inStations = false
			continue
		}

		if !inStations && strings.TrimSpace(line) != "" {
			connection := strings.TrimSpace(line)
			connections = append(connections, connection)

		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error scanning file:", err)
	}

	if !startFound {
		fmt.Println("start station not found")
	}

	if !endFound {
		fmt.Println("end station not found")
	}

	if len(connections) == 0 {
		fmt.Println("no connections found")
	}

	//find routes for trains
	for i := 1; i <= atoi(trainCount); i++ {
		fmt.Printf("\ntrain route %d:\n", i)
		foundRoute := false
		for _, conn := range connections {
			if strings.Contains(conn, startStation) && strings.Contains(conn, endStation) {
				fmt.Println(conn)
				foundRoute = true
			}
		}
		if !foundRoute {
			fmt.Println("no route found", i)
		}
	}
}
func atoi(s string) int {
	num := 0
	for _, c := range s {
		num = num*10 + int(c-'0')
	}
	return num
}

func parseFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inStations := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, "stations:") {
			inStations = true
			continue
		}

		if inStations {
			parts := strings.Split(line, ",")
			if len(parts) >= 3 {
				stationName := parts[0]
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				stations[stationName] = &Station{
					Name:     stationName,
					X:        x,
					Y:        y,
					Adjacent: make(map[string]int),
				}
			}
		}

		if strings.HasPrefix(line, "connections:") {
			inStations = false
			continue
		}

		if !inStations && strings.TrimSpace(line) != "" {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				station1 := strings.TrimSpace(parts[0])
				station2 := strings.TrimSpace(parts[1])
				if _, ok := stations[station1]; ok {
					stations[station1].Adjacent[station2] = 1 // Assume unit weight for simplicity
				}
				if _, ok := stations[station2]; ok {
					stations[station2].Adjacent[station1] = 1 // Bidirectional for undirected graph
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func findRoutes(start, end string) [][]string {
	queue := [][]string{{start}}
	var routes [][]string

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		last := path[len(path)-1]

		if last == end {
			routes = append(routes, path)
		}

		for neighbor := range stations[last].Adjacent {
			if !visited(path, neighbor) {
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return routes
}

func visited(path []string, station string) bool {
	for _, s := range path {
		if s == station {
			return true
		}
	}
	return false
}
