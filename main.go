//main.go
package main

import (
	"fmt"
	"os"
	network "stations/go/network/dijkstra"
	"stations/go/parser"
	"stations/go/pathfinder"

	"strconv"
)

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

	connections, err := parser.ReadMap(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if !isValidStation(connections, startStation) {
		fmt.Fprintln(os.Stderr, "Error: Start station does not exist")
		return
	}

	if !isValidStation(connections, endStation) {
		fmt.Fprintln(os.Stderr, "Error: End station does not exist")
		return
	}


	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station cannot be the same")
		return
	}

	movements := pathfinder.ScheduleTrainMovements(startStation, endStation, connections, numTrains)

	fmt.Print("\nTrain movements from\033[1m ", filePath)
	fmt.Print("\n\033[0m")
	fmt.Print("\033[4m", startStation, "\033[0m to \033[4m", endStation, "\033[0m with \033[4m", numTrains, "\033[0m trains:\n\n")
	for _, move := range movements {
		fmt.Println(move)
	}

	fmt.Printf("\nTotal Movements: %d\n", len(movements))
	fmt.Println("******************************************")
}

func isValidStation(connections network.Connections, station string) bool {
	// Create a map of station names for quick lookup
	stationNames := make(map[string]struct{})
	for _, connection := range connections {
		stationNames[connection.Start.Name] = struct{}{}
		stationNames[connection.End.Name] = struct{}{}
	}
	_, exists := stationNames[station]
	return exists
}
