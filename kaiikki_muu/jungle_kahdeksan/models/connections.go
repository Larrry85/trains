package models

// Connection represents a connection between two stations with a travel time.
type Connection struct {
	Start string
	End   string
	Time  int
}

// Connections is a slice of Connection.
type Connections []Connection
