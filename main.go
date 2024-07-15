package main

import (
	"fmt"
	"os"
	"stations/parser"
	"stations/pathfinder"
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

	connections, err := readMap(filePath)
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

	fmt.Println("Train movements:")
	for _, move := range movements {
		fmt.Println(move)
	}

	fmt.Printf("\nTotal number of movements: %d\n", len(movements))
}

func readMap(filePath string) (pathfinder.Connections, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parser.ParseConnections(file)
}

func isValidStation(connections pathfinder.Connections, station string) bool {
	for _, connection := range connections {
		if connection.Start == station || connection.End == station {
			return true
		}
	}
	return false
}
