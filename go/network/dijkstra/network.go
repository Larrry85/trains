package network

// Station represents a station with an X, Y coordinate.
type Station struct {
	Name string
	X, Y int
}

// Item represents an element in the priority queue with a value, priority, and index.
type Item struct {
	Value    string
	Priority int
	index    int
	Station  string
	Cost     int
	Parent   *Item
	Time     int
	TrainID  int
}

// PriorityQueue implements a priority queue using a slice of Items.
type PriorityQueue []*Item

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less returns true if the priority of item i is less than the priority of item j.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

// Swap swaps the elements with indexes i and j in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item x to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Connection represents a connection between two stations with a travel time.
type Connection struct {
	Start Station
	End   Station
	Time  int
}

// Connections is a slice of Connection.
type Connections []Connection

// Train represents a train with an ID and color.
type Train struct {
	ID    int
	Color string
}
