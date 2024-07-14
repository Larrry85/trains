package parser

import (
	"bufio"
	"io"
	"strings"
)

type Connection struct {
	Start string
	End   string
}

type Connections []Connection

func ParseConnections(r io.Reader) (Connections, error) {
	scanner := bufio.NewScanner(r)
	connections := Connections{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			continue
		}
		connections = append(connections, Connection{
			Start: parts[0],
			End:   parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}
