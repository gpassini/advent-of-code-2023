package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

const (
	left  Direction = 'L'
	right Direction = 'R'
)

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))

	var directions Directions
	r.Scan()
	for _, c := range r.Text() {
		directions = append(directions, NewDirection(c))
	}
	fmt.Println("Directions:", directions)

	r.Scan() // empty line

	labelToNode := make(map[string]*Node)
	destsToLabel := make(map[string]string)
	replaceLabel := make(map[string]string)
	var firstNode *Node
	for r.Scan() {
		node := parseNode(r.Text())
		if node.label == "AAA" {
			firstNode = node
		}
		dests := fmt.Sprintf("%s%s", node.leftLabel, node.rightLabel)
		if label, ok := destsToLabel[dests]; ok {
			fmt.Printf("Replaced %s with %s\n", node.label, label)
			replaceLabel[node.label] = label
			labelToNode[label] = node
		} else {
			destsToLabel[dests] = node.label
			labelToNode[node.label] = node
		}
	}
	if firstNode == nil {
		panic("AAA node not found")
	}

	var stepsCount int
	directionIter := directions.Iter()
	node := firstNode
	for {
		fmt.Println("Node:", node)
		if node.label == "ZZZ" {
			break
		}

		var steps int
		node, steps = node.Move(directionIter, labelToNode, replaceLabel)
		stepsCount += steps
	}

	println(stepsCount)
}

type Node struct {
	label      string
	leftLabel  string
	rightLabel string
	visited    map[string]*Node
}

func (n *Node) Visit(direction Direction, labelToNode map[string]*Node, replaceLabel map[string]string) *Node {
	switch direction {
	case left:
		leftLabel := n.leftLabel
		if label, ok := replaceLabel[leftLabel]; ok {
			leftLabel = label
		}
		if leftNode, ok := labelToNode[leftLabel]; !ok {
			panic(fmt.Sprintf("node not found with label: %s", leftLabel))
		} else {
			n.visited[direction.String()] = leftNode
			return leftNode
		}
	case right:
		rightLabel := n.rightLabel
		if label, ok := replaceLabel[rightLabel]; ok {
			rightLabel = label
		}
		if rightNode, ok := labelToNode[rightLabel]; !ok {
			panic(fmt.Sprintf("node not found with label: %s", rightLabel))
		} else {
			n.visited[direction.String()] = rightNode
			return rightNode
		}
	default:
		panic(fmt.Sprintf("invalid direction: %c", direction))
	}
}

func (n *Node) Move(directionsIter DirectionsIter, labelToNode map[string]*Node, replaceLabel map[string]string) (*Node, int) {
	var steps int
	var directions Directions
	furthestNode := n
	for {
		if furthestNode.label == "ZZZ" {
			return furthestNode, steps
		}
		latestDirection := directionsIter()
		steps++
		directions = append(directions, latestDirection)
		directionsStr := directions.String()
		if maybeFurthestNode, ok := n.visited[directionsStr]; ok {
			furthestNode = maybeFurthestNode
		} else {
			furthestNode := furthestNode.Visit(latestDirection, labelToNode, replaceLabel)
			n.visited[directionsStr] = furthestNode
			fmt.Println("Evaluated:", directionsStr)
			return furthestNode, steps
		}
	}
}

func (n Node) String() string {
	return fmt.Sprintf("%s = (%s, %s)", n.label, n.leftLabel, n.rightLabel)
}

type Direction rune

func NewDirection(c rune) Direction {
	switch c {
	case 'L', 'R':
		return Direction(c)
	}
	panic(fmt.Sprintf("invalid direction: %c", c))
}

func (d Direction) String() string {
	return string(d)
}

type Directions []Direction

func (d Directions) String() string {
	var sb strings.Builder
	for _, c := range d {
		sb.WriteRune(rune(c))
	}
	return sb.String()
}

type DirectionsIter func() Direction

func (d Directions) Iter() DirectionsIter {
	var idx int
	return func() Direction {
		if idx >= len(d) {
			idx = 0
		}
		next := d[idx]
		idx++
		return next
	}
}

func parseNode(s string) *Node {
	labelAndDirections := strings.Split(s, " = ")
	if len(labelAndDirections) != 2 {
		panic(fmt.Sprintf("unexpected split size (%d) for string: %s", len(labelAndDirections), s))
	}
	label := labelAndDirections[0]
	directions := labelAndDirections[1][1 : len(labelAndDirections[1])-1]
	leftRightLabels := strings.Split(directions, ", ")
	if len(leftRightLabels) != 2 {
		panic(fmt.Sprintf("unexpected split size (%d) for directions: %s", len(leftRightLabels), directions))
	}

	leftLabel := leftRightLabels[0]
	rightLabel := leftRightLabels[1]
	return &Node{
		label:      label,
		leftLabel:  leftLabel,
		rightLabel: rightLabel,
		visited:    make(map[string]*Node),
	}
}
