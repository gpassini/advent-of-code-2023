package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

const (
	N_S    Pipe = '|'
	E_W    Pipe = '-'
	N_E    Pipe = 'L'
	N_W    Pipe = 'J'
	S_E    Pipe = 'F'
	S_W    Pipe = '7'
	Ground Pipe = '.'
	Start  Pipe = 'S'
)

const (
	Up Direction = iota
	Down
	Right
	Left
)

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))

	var pipes Pipes
	var startingPoint Point

	var lineIdx int
	for r.Scan() {
		var pipesLine []Pipe
		for columnIdx, c := range r.Text() {
			pipesLine = append(pipesLine, Pipe(c))
			if Pipe(c) == Start {
				startingPoint = Point{x: lineIdx, y: columnIdx}
			}
		}
		pipes = append(pipes, pipesLine)
		lineIdx++
	}

	var steps int
	pos := startingPoint
	lastDir := Up
	for {
		fmt.Println("Pos:", pos)
		pos, lastDir = move(pipes, pos, lastDir)
		steps++
		if pipes[pos.x][pos.y] == Start {
			break
		}
	}

	println(steps / 2)
}

type Pipe rune

func (p Pipe) Connects(dir Direction) bool {
	return slices.Contains(p.Exits(), dir)
}

func (p Pipe) Exits() []Direction {
	switch p {
	case N_S:
		return []Direction{Up, Down}
	case E_W:
		return []Direction{Right, Left}
	case N_E:
		return []Direction{Up, Right}
	case N_W:
		return []Direction{Up, Left}
	case S_E:
		return []Direction{Down, Right}
	case S_W:
		return []Direction{Down, Left}
	case Ground:
		return []Direction{}
	case Start:
		return []Direction{Up, Down, Right, Left}
	default:
		panic(fmt.Sprintf("invalid pipe: %c", p))
	}
}

type Pipes [][]Pipe

type Point struct {
	x, y int
}

func (p Point) Apply(d Direction) Point {
	switch d {
	case Up:
		return Point{x: p.x - 1, y: p.y}
	case Down:
		return Point{x: p.x + 1, y: p.y}
	case Right:
		return Point{x: p.x, y: p.y + 1}
	case Left:
		return Point{x: p.x, y: p.y - 1}
	default:
		panic(fmt.Sprintf("invalid direction: %d", d))
	}
}

type Direction int

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Right:
		return Left
	case Left:
		return Right
	default:
		panic(fmt.Sprintf("invalid direction: %d", d))
	}
}

func move(pipes Pipes, position Point, lastDirection Direction) (Point, Direction) {
	x, y := position.x, position.y
	currentPipe := pipes[x][y]

	if currentPipe == Start {
		// look for any connecting pipe
		if dir := Up; x > 0 && pipes[x-1][y].Connects(dir.Opposite()) {
			return position.Apply(dir), dir
		}
		if dir := Down; x < len(pipes)-1 && pipes[x+1][y].Connects(dir.Opposite()) {
			return position.Apply(dir), dir
		}
		if dir := Left; y > 0 && pipes[x][y-1].Connects(dir.Opposite()) {
			return position.Apply(dir), dir
		}
		if dir := Right; y < len(pipes[x])-1 && pipes[x][y+1].Connects(dir.Opposite()) {
			return position.Apply(dir), dir
		}
	}

	availableDirs := currentPipe.Exits()
	for _, dir := range availableDirs {
		if dir != lastDirection.Opposite() {
			return position.Apply(dir), dir
		}
	}

	panic(fmt.Sprintf("No exit found (%v, %v)", position, lastDirection))
}
