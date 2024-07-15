/*package parser

import (
	"bufio"
	"io"
	"stations/pathfinder"
	"strings"
)

func ParseConnections(r io.Reader) (pathfinder.Connections, error) {
	scanner := bufio.NewScanner(r)
	connections := pathfinder.Connections{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			continue
		}
		connections = append(connections, pathfinder.Connection{
			Start: parts[0],
			End:   parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}
*/