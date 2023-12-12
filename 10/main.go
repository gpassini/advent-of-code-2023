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

	In  Pipe = 'I'
	Out Pipe = 'O'
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

	allPipes = []Pipe{N_S, E_W, N_E, N_W, S_E, S_W}
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))

	var pipes Pipes
	var isLoopPartMask [][]bool
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
		isLoopPartMask = append(isLoopPartMask, make([]bool, len(pipesLine)))
		lineIdx++
	}

	var steps int
	pos := startingPoint
	lastDir := Up
	firstDir := Direction('@')
	for {
		fmt.Println("Pos:", pos)
		isLoopPartMask[pos.x][pos.y] = true
		pos, lastDir = move(pipes, pos, lastDir)
		if firstDir == Direction('@') {
			firstDir = lastDir
		}
		steps++
		if pipes[pos.x][pos.y] == Start {
			break
		}
	}

	println(steps / 2)

	startingPipe := FromDirections(firstDir, lastDir.Opposite())
	pipes[startingPoint.x][startingPoint.y] = startingPipe

	var res int
	for i, pipesLine := range pipes {
		inCount := paintInOut(pipesLine, isLoopPartMask[i])
		fmt.Printf("Line %d has %d inside\n", i, inCount)
		res += inCount
	}

	fmt.Println(pipes)

	println(res)
}

type Pipe rune

func FromDirections(d1, d2 Direction) Pipe {
	if d1 == d2 {
		panic("same directions given")
	}
	for _, p := range allPipes {
		exits := p.Exits()
		if slices.Contains(exits, d1) && slices.Contains(exits, d2) {
			return p
		}
	}
	panic("could not find pipe for given directions")
}

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

func (p Pipe) String() string {
	return string(p)
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

func (p Pipes) String() string {
	var sb strings.Builder
	for _, l := range p {
		for _, c := range l {
			sb.WriteRune(rune(c))
		}
		sb.WriteString("\n")
	}
	return sb.String()
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

func paintInOut(pipes []Pipe, isLoopMask []bool) (inCount int) {
	var inside bool
	var firstWallPipe Pipe
	for i, p := range pipes {
		if isLoopMask[i] {
			switch p {
			case N_S:
				inside = !inside
			case N_E, S_E:
				// entering wall
				firstWallPipe = p
			case N_W:
				// exiting wall
				if firstWallPipe == S_E {
					inside = !inside
				}
			case S_W:
				// exiting wall
				if firstWallPipe == N_E {
					inside = !inside
				}
			}
		} else {
			if inside {
				pipes[i] = In
				inCount++
			} else {
				pipes[i] = Out
			}
		}
	}
	return inCount
}
