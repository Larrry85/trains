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

### Invalid maps

There ae maps that contain errors, for example:

```
go run . maps/errors/10no-start-station_london.txt waterloo victoria 4
```


### OIKEAT


T1-victoria T2-euston

T1-st_pancras T2-st_pancras T3-victoria T4-euston

T3-st_pancras T4-st_pancras

Movements: 

***********


T1-apple_avenue

T1-orange_junction T2-apple_avenue

T1-space_port T2-orange_junction T3-apple_avenue

T2-space_port T3-orange_junction T4-apple_avenue

T3-space_port T4-orange_junction

T4-space_port

Movements: 6

***********

T1-grasslands T2-farms T3-green_belt

T1-suburbs T2-downtown T3-village T4-grasslands T5-farms T6-green_belt

T1-clouds T2-metropolis T3-mountain T4-suburbs T5-downtown T6-village T7-grasslands T8-farms T9-green_belt

T1-wetlands T2-industrial T3-treetop T4-clouds T5-metropolis T6-mountain T7-suburbs T8-downtown T9-village T10-grasslands

T1-desert T2-desert T3-desert T4-wetlands T5-industrial T6-treetop T7-clouds T8-metropolis T9-mountain T10-suburbs

T4-desert T5-desert T6-desert T7-wetlands T8-industrial T9-treetop T10-clouds

T7-desert T8-desert T9-desert T10-wetlands

T10-desert

Movements: 8

***********

T1-terminus T2-near

T2-far T3-terminus T4-near

T2-terminus T4-far T5-terminus T6-near

T4-terminus T6-far T7-terminus T8-near

T6-terminus T8-far T9-terminus T10-near

T8-terminus T10-far T11-terminus T12-near

T10-terminus T12-far T13-terminus T14-near

T12-terminus T14-far T15-terminus T16-near

T14-terminus T16-far T17-terminus T18-near

T16-terminus T18-far T19-terminus

T18-terminus T20-terminus

Movements: 11

***********

T1-three

T1-one T2-three

T1-four T2-one T3-three

T2-four T3-one T4-three

T3-four T4-one

T4-four

Movements: 6

***********

T1-verdi T3-handel

T1-part T2-verdi T3-mozart T5-handel

T2-part T3-part T4-verdi T5-mozart T7-handel

T4-part T5-part T6-verdi T7-mozart T9-handel

T6-part T7-part T8-verdi T9-mozart

T8-part T9-part

Movements: 6

***********

T1-10 T4-13 T6-00

T1-11 T2-10 T4-14 T5-13 T6-01

T1-12 T2-11 T3-10 T4-15 T5-14 T6-02 T9-13

T1-large T2-12 T3-11 T4-21 T5-15 T6-03 T7-10 T9-14

T2-large T3-12 T4-22 T5-21 T6-04 T7-11 T8-10 T9-15

T3-large T4-large T5-22 T6-05 T7-12 T8-11 T9-21

T5-large T6-large T7-large T8-12 T9-22

T8-large T9-large

Movements: 8

***********


## Coders

Laura Levistö - Jonathan Dahl - 6/24