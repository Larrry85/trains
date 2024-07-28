package astar

type PriorityQueue []*Node

type Node struct {
	Station  string
	Cost     int
	Priority int
	Parent   *Node
	Time     int
}

type Station struct {
	Name string
	X    int
	Y    int
}

type Graph struct {
	Stations    map[string]Station
	Connections map[string][]string
}

// Implement heap.Interface for PriorityQueue
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*Node)
	*pq = append(*pq, node)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

// Train represents a train with an ID and color.
type Train struct {
	ID    int
	Color string
}

// Queue implementation for station queues with string support
type StringQueue struct {
	items []string
}

func NewStringQueue() *StringQueue {
	return &StringQueue{items: []string{}}
}

func (q *StringQueue) Push(item string) {
	q.items = append(q.items, item)
}

func (q *StringQueue) Pop() {
	if len(q.items) == 0 {
		return
	}
	q.items = q.items[1:]
}

func (q *StringQueue) Front() string {
	if len(q.items) == 0 {
		return ""
	}
	return q.items[0]
}

func (q *StringQueue) Remove(val string) {
	for i := 0; i < len(q.items); i++ {
		if q.items[i] == val {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return
		}
	}
}
