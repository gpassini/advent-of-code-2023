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

	var aNodes []*Node
	for r.Scan() {
		node := parseNode(r.Text())
		if strings.HasSuffix(node.label, "A") {
			aNodes = append(aNodes, node)
		}
		labelToNode[node.label] = node
	}
	if len(aNodes) == 0 {
		panic("A nodes not found")
	}

	periodEndSteps := make([]int, len(aNodes))
	for i, node := range aNodes {
		iter := directions.Iter()
		node, directions := node.MoveIter(iter, labelToNode)
		if !strings.HasSuffix(node.label, "Z") {
			panic("expected end node")
		}
		periodEndSteps[i] = len(directions)
		fmt.Println("Period:", i, periodEndSteps[i])
	}

	println(lcm(periodEndSteps[0], periodEndSteps[1], periodEndSteps[2:]...))
}

type Node struct {
	label      string
	leftLabel  string
	rightLabel string
	visited    map[string]*Node
}

func (n *Node) Visit(direction Direction, labelToNode map[string]*Node) *Node {
	switch direction {
	case left:
		leftLabel := n.leftLabel
		if leftNode, ok := labelToNode[leftLabel]; !ok {
			panic(fmt.Sprintf("node not found with label: %s", leftLabel))
		} else {
			n.visited[direction.String()] = leftNode
			return leftNode
		}
	case right:
		rightLabel := n.rightLabel
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

func (n *Node) MoveIter(directionsIter DirectionsIter, labelToNode map[string]*Node) (*Node, Directions) {
	var directions Directions
	furthestNode := n
	for {
		if len(directions) > 0 && strings.HasSuffix(furthestNode.label, "Z") {
			return furthestNode, directions
		}
		latestDirection := directionsIter()
		directions = append(directions, latestDirection)
		directionsStr := directions.String()
		if maybeFurthestNode, ok := n.visited[directionsStr]; ok {
			furthestNode = maybeFurthestNode
		} else {
			furthestNode = furthestNode.Visit(latestDirection, labelToNode)
			n.visited[directionsStr] = furthestNode
		}
	}
}

func (n *Node) Move(givenDirections Directions, labelToNode map[string]*Node) *Node {
	furthestNode := n
	var directions Directions
	for _, latestDirection := range givenDirections {
		directions = append(directions, latestDirection)
		directionsStr := directions.String()
		if maybeFurthestNode, ok := n.visited[directionsStr]; ok {
			furthestNode = maybeFurthestNode
		} else {
			furthestNode = furthestNode.Visit(latestDirection, labelToNode)
			n.visited[directionsStr] = furthestNode
		}
	}
	return furthestNode
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

func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
