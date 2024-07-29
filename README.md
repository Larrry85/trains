# Stations

This Go program simulates the movement of multiple trains on a rail network, from a starting station to an ending station, based on the shortest paths between the stations.

## Before starting the program

- You have to have Golang installed.

## How to use

Usage: go run . [path to file containing network map] [start station] [end station] [number of trains]

### Valid maps

There are valid train routes, for example:

 ```
go run . maps/01london.txt waterloo st_pancras 2
 ```
Text file of correct test cases is found in the maps directory (tests.txt). 

### Invalid maps

There are maps that contain errors, for example:

```
go run . maps/errors/10no-start-station_london.txt waterloo victoria 4
```
Test text file of incorrect cases is found in the maps/errors directory (tests_errors.txt).

## Directory tree and explanations of the GO files

```
stations/
│
├── go/
│   ├── A/
│   │   └── A.go
│   ├── network/
│   │   ├── astar/
│   │   │   └── Anetwork.go
│   │   └── dijkstra/
│   │       └── network.go
│   ├── parser/
│   │   └── parser.go
│   └── pathfinder/
│   │   └── pathfinder.go
├── maps/
│   ├── errors/
│   │   └── tests_errors.txt  
│   └── tests.txt              
├── go.mod 
├── main.go
└── README.md
```               

main.go:
- Validates input arguments.
- Reads and parses the train map text file.
- Uses ScheduleTrainMovements() from pathfinder.go.
- Prints the total movements.

A.go:
- Handles the most trickiest train map, 07small.txt.
- Reads train map text file and validates the content.
- findDistinctPaths() identifies distinct paths between a start and end station.
- aStarPathfinding() finds the optimal path using the A* algorithm.
- distributeTrainsAcrossPaths() simulates train movements, identifies the shortest, second shortest, and longest paths and defines the amount of paths. Prints the total movements.

Anetwork.go:
- Data structs for A* pathfinding algorithm: Node(a node in the graph, a step in the potential path, the current state in the search process), Station(one station), and Graph(network  of stations and connections), Train(trains in the simulation, id and color), StringQueue(station names).
- PriorityQueue to manage nodes based on their priorities.

network.go:
- Data structs: Station(one station), Item(element in the priority queue), Connection(connection between two stations), Train(trains in the simulation, id and color)
- PriorityQueue for efficient pathfinding and scheduling.

parser.go:
- Reads train map text file and validates the content.
- Constructs a network representation for use in the Dijkstra pathfinding algorithm.

pathfinder.go:
- FindShortestPath() finds the shortest path from start to end using Dijkstra algorithm.
- FindAllPaths() finds all possible paths from start to end station.
- Heurestic() calculates the Euclidean distance between two stations, (heuristic in pathfinding).
- buildAdjacencyList() converts a list of connections into an adjacency list representation, mapping each station to its neighboring stations with travel times.
- ScheduleTrainMovements() simulates train movements from start to end station using, for example, Dijkstra algorithm.

## Coders

Laura Levistö - Jonathan Dahl - 7/24