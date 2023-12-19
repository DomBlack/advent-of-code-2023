package day08

import (
	"io"

	"github.com/DomBlack/advent-of-code-2023/pkg/maths"
	"github.com/DomBlack/advent-of-code-2023/pkg/runner"
	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var Day08 = runner.NewDay(8, parseMap, part1, part2)

func part1(_ *runner.Context, _ zerolog.Logger, input Map) (answer int, err error) {
	// Follow the instructions
	node := input.Root
	steps := 0
	for {
		if node.Name == "ZZZ" {
			break
		}

		switch input.Instructions[steps%len(input.Instructions)] {
		case Left:
			node = node.Left
		case Right:
			node = node.Right
		}

		steps++
	}

	return steps, nil
}

func part2(_ *runner.Context, _ zerolog.Logger, input Map) (answer int, err error) {
	lengths := make([]int, len(input.NodesEndingWithA))

	// Calculate the length of each of the parallel paths
	for i, node := range input.NodesEndingWithA {
		steps := 0
		for {
			// Check if all the position ends with Z
			if node.Name[2] == 'Z' {
				break
			}

			// Move the position
			switch input.Instructions[steps%len(input.Instructions)] {
			case Left:
				node = node.Left
			case Right:
				node = node.Right
			}

			steps++
		}

		lengths[i] = steps
	}

	// Now find the lowest common multiple of all the lengths
	return maths.LCM(lengths), nil
}

type Instruction uint8

const (
	Left Instruction = iota
	Right
)

type Map struct {
	Instructions []Instruction

	Root             *Node
	NodesEndingWithA []*Node
}

type Node struct {
	Name  string
	Left  *Node
	Right *Node
}

func parseMap(input []byte) (rtn Map, err error) {
	lines := stream.LinesFrom(input)

	// First read the instructions
	instructions, err := lines.Next()
	if err != nil {
		return Map{}, errors.Wrap(err, "failed to read instructions")
	}

	rtn.Instructions = make([]Instruction, len(instructions))
	for i, instruction := range instructions {
		switch instruction {
		case 'L':
			rtn.Instructions[i] = Left
		case 'R':
			rtn.Instructions[i] = Right
		default:
			return Map{}, errors.Errorf("invalid instruction %c", instruction)
		}
	}

	// Read the blank line
	line, err := lines.Next()
	if err != nil {
		return Map{}, errors.Wrap(err, "failed to read blank line")
	}
	if line != "" {
		return Map{}, errors.Errorf("expected blank line, got %q", line)
	}

	// Create a map of all the nodes
	nodes := make(map[string]*Node)
	getNode := func(name string) *Node {
		if node, ok := nodes[name]; ok {
			return node
		}

		node := &Node{Name: name}
		nodes[name] = node

		return node
	}

	// Loop over the lines constructing the map
	for {
		line, err := lines.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Map{}, errors.Wrap(err, "failed to read line")
		}

		// AAA = (BBB, BBB)
		name := line[:3]
		leftName := line[7:10]
		rightName := line[12:15]

		node := getNode(name)
		node.Left = getNode(leftName)
		node.Right = getNode(rightName)

		if name == "AAA" {
			rtn.Root = node
		}

		if name[2] == 'A' {
			rtn.NodesEndingWithA = append(rtn.NodesEndingWithA, node)
		}
	}

	// Ensure all nodes are connected
	for _, node := range nodes {
		if node.Left == nil || node.Right == nil {
			return Map{}, errors.Errorf("node %s is not connected", node.Name)
		}
	}

	return rtn, nil
}
