// main.go
package main

import (
	"fmt"
	"os"
	"stations/train" // Assuming train package is located in stations/train directory
	"strconv"
)

func main() {
	// Check if the number of command-line arguments is correct
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Error: Incorrect number of arguments")
		return
	}

	// Extract command-line arguments
	filePath := os.Args[1]                     // Path to the file containing station connections
	startStation := os.Args[2]                 // Starting station
	endStation := os.Args[3]                   // Ending station
	numTrains, err := strconv.Atoi(os.Args[4]) // Number of trains, converted from string to int
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of trains must be a positive integer")
		return
	}

	// Read the station connections from the specified file
	connections, err := readMap(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Check if the start station exists in the connections
	if !isValidStation(connections, startStation) {
		fmt.Fprintln(os.Stderr, "Error: Start station does not exist")
		return
	}

	// Check if the end station exists in the connections
	if !isValidStation(connections, endStation) {
		fmt.Fprintln(os.Stderr, "Error: End station does not exist")
		return
	}

	// Check if start station and end station are the same
	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station cannot be the same")
		return
	}

	// Calculate train movements using the ScheduleTrainMovements function from the train package
	movements := train.ScheduleTrainMovements(startStation, endStation, connections, numTrains)

	// Output the train movements
	fmt.Println("Train movements:")
	for _, move := range movements {
		fmt.Println(move)
	}

	// Output the total number of movements
	fmt.Printf("\nTotal number of movements: %d\n", len(movements))
}

// Function to read station connections from a file
func readMap(filePath string) (train.Connections, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return train.ParseConnections(file)
}

// Function to check if a station exists in the connections
func isValidStation(connections train.Connections, station string) bool {
	for _, connection := range connections {
		if connection.Start == station || connection.End == station {
			return true
		}
	}
	return false
}
