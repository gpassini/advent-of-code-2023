package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

var (
	//go:embed input.txt
	input string
)

func main() {
	contraption := Contraption{
		energized: make(map[Coordinate][]Direction),
	}
	r := bufio.NewScanner(strings.NewReader(input))
	for r.Scan() {
		var tilesLine []Tile
		for _, c := range r.Text() {
			tilesLine = append(tilesLine, Tile(c))
		}
		contraption.tiles = append(contraption.tiles, tilesLine)
	}

	tilesToVisit := []TileToVisit{{
		rowIdx:    0,
		columnIdx: 0,
		dir:       right,
	}}
	for len(tilesToVisit) > 0 {
		nextTile := tilesToVisit[0]
		newTilesToVisit := contraption.Visit(nextTile.rowIdx, nextTile.columnIdx, nextTile.dir)
		tilesToVisit = append(tilesToVisit[1:], newTilesToVisit...)
	}

	println(len(contraption.energized))
}

type Direction uint

const (
	up Direction = iota
	down
	left
	right
)

func (d Direction) Move(rowIdx, columnIdx int) (int, int) {
	switch d {
	case up:
		rowIdx--
	case down:
		rowIdx++
	case left:
		columnIdx--
	case right:
		columnIdx++
	default:
		panic(fmt.Sprintf("invalid direction: %d", d))
	}
	return rowIdx, columnIdx
}

type Tile rune

const (
	emptySpace         Tile = '.'
	forwardMirror      Tile = '/'
	backwardMirrot     Tile = '\\'
	verticalSplitter   Tile = '|'
	horizontalSplitter Tile = '-'
)

func (t Tile) Reflect(dir Direction) []Direction {
	switch t {
	case emptySpace:
		return []Direction{dir}
	case forwardMirror:
		switch dir {
		case up:
			return []Direction{right}
		case down:
			return []Direction{left}
		case left:
			return []Direction{down}
		case right:
			return []Direction{up}
		}
	case backwardMirrot:
		switch dir {
		case up:
			return []Direction{left}
		case down:
			return []Direction{right}
		case left:
			return []Direction{up}
		case right:
			return []Direction{down}
		}
	case verticalSplitter:
		switch dir {
		case up, down:
			return []Direction{dir}
		case left, right:
			return []Direction{up, down}
		}
	case horizontalSplitter:
		switch dir {
		case up, down:
			return []Direction{left, right}
		case left, right:
			return []Direction{dir}
		}
	}
	panic(fmt.Sprintf("Tile %d reflection not handled for direction %d", t, dir))
}

type Contraption struct {
	tiles     [][]Tile
	energized map[Coordinate][]Direction
}

func (contraption Contraption) Visit(rowIdx, columnIdx int, dir Direction) []TileToVisit {
	if rowIdx < 0 || rowIdx >= len(contraption.tiles) || columnIdx < 0 || columnIdx >= len(contraption.tiles[0]) {
		return nil
	}

	if dirs, ok := contraption.energized[Coordinate{rowIdx, columnIdx}]; ok {
		for _, d := range dirs {
			if d == dir {
				// already visited this tile with this direction
				return nil
			}
		}
		// add direction to the list of visited directions for this tile
		contraption.energized[Coordinate{rowIdx, columnIdx}] = append(contraption.energized[Coordinate{rowIdx, columnIdx}], dir)
	} else {
		contraption.energized[Coordinate{rowIdx, columnIdx}] = []Direction{dir}
	}

	currTile := contraption.tiles[rowIdx][columnIdx]
	newDirs := currTile.Reflect(dir)
	var tilesToVisit []TileToVisit
	for _, newDir := range newDirs {
		newRowIdx, newColumnIdx := newDir.Move(rowIdx, columnIdx)
		tilesToVisit = append(tilesToVisit, TileToVisit{
			rowIdx:    newRowIdx,
			columnIdx: newColumnIdx,
			dir:       newDir,
		})
	}
	return tilesToVisit
}

type Coordinate struct{ rowIdx, columnIdx int }

type TileToVisit struct {
	rowIdx, columnIdx int
	dir               Direction
}
