# Stations

This Go program simulates the movement of multiple trains on a rail network, from a starting station to an ending station, based on the shortest paths between the stations

## Before starting the program

- You have to have Golang installed.

## How to use

Usage: go run . [path to file containing network map] [start station] [end station] [number of trains]

### Valid maps

There are valid train routes, for example:

 ```
go run . maps/01london.txt waterloo st_pancras 2
 ```
Test text file of correct cases is found in the maps directory. 

### Invalid maps

There are maps that contain errors, for example:

```
go run . maps/errors/10no-start-station_london.txt waterloo victoria 4
```
Test text file of incorrect cases is found in the maps/errors directory.

## Coders

Laura Levistö - Jonathan Dahl - 7/24