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
	dish := parseDish()
	const nCycles = 1_000_000_000
	knownPositions := make(map[string]int)
	var firstRepNumber int
	var repPeriod int
	for i := 0; i < nCycles; i++ {
		if knownPosIdx, ok := knownPositions[dish.String()]; ok {
			fmt.Println("Position cycle:", knownPosIdx, i)
			firstRepNumber = knownPosIdx
			repPeriod = i - knownPosIdx
			break
		} else {
			knownPositions[dish.String()] = i
		}

		if i%1_000_000 == 0 {
			fmt.Println("Tilted", i, "times:\n", dish)
		}
		dish.Tilt(north)
		dish.Tilt(west)
		dish.Tilt(south)
		dish.Tilt(east)

		// if i <= 2 {
		// 	fmt.Println("After", i+1, "cycle(s):\n", dish)
		// }
	}
	neededCycles := (nCycles - firstRepNumber) % repPeriod
	fmt.Println("Needed cycles:", neededCycles)
	for i := 0; i < neededCycles; i++ {
		dish.Tilt(north)
		dish.Tilt(west)
		dish.Tilt(south)
		dish.Tilt(east)
	}
	fmt.Println("Result:", dish.Evaluate(north))
}

func parseDish() Dish {
	var dish Dish
	r := bufio.NewScanner(strings.NewReader(input))
	for r.Scan() {
		var tilesRow []Tile
		for _, c := range r.Text() {
			tilesRow = append(tilesRow, Tile(c))
		}
		dish = append(dish, tilesRow)
	}
	return dish
}

type Direction uint

const (
	north Direction = iota
	south
	east
	west
)

func (d Direction) MoveBack(rowIdx, columnIdx int) (int, int) {
	switch d {
	case north:
		rowIdx++
	case south:
		rowIdx--
	case west:
		columnIdx++
	case east:
		columnIdx--
	default:
		panic(fmt.Sprintf("unknown direction: %d", d))
	}
	return rowIdx, columnIdx
}

type Tile rune

const (
	roundedRock    Tile = 'O'
	cubeShapedRock Tile = '#'
	emptySpace     Tile = '.'
)

type Dish [][]Tile

func (d *Dish) Tilt(dir Direction) {
	rowLenght := len((*d)[0])
	columnLength := len(*d)
	var tiltDirectionWidth int
	var axisDirectionIdx func(int, int) int
	// var tiltDirectionIdx func(int, int) int
	var updateTiltDirectionIdx func(int, int) (int, int)
	var updateAxisDirectionIdx func(int, int) (int, int)
	var startRowIdx, startColumnIdx int
	switch dir {
	case north:
		tiltDirectionWidth = rowLenght
		axisDirectionIdx = func(_, columnIdx int) int { return columnIdx }
		// tiltDirectionIdx = func(rowIdx, _ int) int { return rowIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx + 1, 0 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx, columnIdx + 1 }
		startRowIdx = 0
		startColumnIdx = 0
	case south:
		tiltDirectionWidth = rowLenght
		axisDirectionIdx = func(_, columnIdx int) int { return columnIdx }
		// tiltDirectionIdx = func(rowIdx, _ int) int { return rowIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx - 1, 0 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx, columnIdx + 1 }
		startRowIdx = columnLength - 1
		startColumnIdx = 0
	case west:
		tiltDirectionWidth = columnLength
		axisDirectionIdx = func(rowIdx, _ int) int { return rowIdx }
		// tiltDirectionIdx = func(_, columnIdx int) int { return columnIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return 0, columnIdx + 1 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx + 1, columnIdx }
		startRowIdx = 0
		startColumnIdx = 0
	case east:
		tiltDirectionWidth = columnLength
		axisDirectionIdx = func(rowIdx, _ int) int { return rowIdx }
		// tiltDirectionIdx = func(_, columnIdx int) int { return columnIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return 0, columnIdx - 1 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx + 1, columnIdx }
		startRowIdx = 0
		startColumnIdx = rowLenght - 1
	default:
		panic(fmt.Sprintf("unknown dir: %d", dir))
	}

	// indexes of the furtherst empty space in that direction
	furthestEmptySpaces := make([]*Pos, tiltDirectionWidth)

	rowIdx, columnIdx := startRowIdx, startColumnIdx
	for {
		for {
			tile := (*d)[rowIdx][columnIdx]
			switch tile {
			case roundedRock:
				furthestEmptySpace := furthestEmptySpaces[axisDirectionIdx(rowIdx, columnIdx)]
				if furthestEmptySpace != nil {
					// fmt.Printf("Rolling rock at (%d, %d)\n", rowIdx, columnIdx)
					(*d)[furthestEmptySpace.rowIdx][furthestEmptySpace.columnIdx] = roundedRock
					(*d)[rowIdx][columnIdx] = emptySpace
					nextEmptySpaceRowIdx, nextEmptySpaceColumnIdx := dir.MoveBack(furthestEmptySpace.rowIdx, furthestEmptySpace.columnIdx)
					furthestEmptySpaces[axisDirectionIdx(rowIdx, columnIdx)] = &Pos{nextEmptySpaceRowIdx, nextEmptySpaceColumnIdx}
					// fmt.Printf("Next empty space at (%d, %d)\n", nextEmptySpaceRowIdx, nextEmptySpaceColumnIdx)
				}
			case cubeShapedRock:
				furthestEmptySpaces[axisDirectionIdx(rowIdx, columnIdx)] = nil
			case emptySpace:
				furthestEmptySpace := furthestEmptySpaces[axisDirectionIdx(rowIdx, columnIdx)]
				if furthestEmptySpace == nil {
					furthestEmptySpaces[axisDirectionIdx(rowIdx, columnIdx)] = &Pos{rowIdx, columnIdx}
				}
			}

			rowIdx, columnIdx = updateAxisDirectionIdx(rowIdx, columnIdx)
			if rowIdx < 0 || rowIdx >= columnLength {
				break
			}
			if columnIdx < 0 || columnIdx >= rowLenght {
				break
			}
		}
		rowIdx, columnIdx = updateTiltDirectionIdx(rowIdx, columnIdx)
		if rowIdx < 0 || rowIdx >= columnLength {
			break
		}
		if columnIdx < 0 || columnIdx >= rowLenght {
			break
		}
	}
}

func (d Dish) Evaluate(dir Direction) int {
	rowLenght := len(d[0])
	columnLength := len(d)
	var axisDirectionWidth int
	var tiltDirectionIdx func(int, int) int
	var updateTiltDirectionIdx func(int, int) (int, int)
	var updateAxisDirectionIdx func(int, int) (int, int)
	var startRowIdx, startColumnIdx int
	switch dir {
	case north:
		axisDirectionWidth = columnLength
		tiltDirectionIdx = func(rowIdx, _ int) int { return rowIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx + 1, 0 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx, columnIdx + 1 }
		startRowIdx = 0
		startColumnIdx = 0
	case south:
		axisDirectionWidth = columnLength
		tiltDirectionIdx = func(rowIdx, _ int) int { return rowIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx - 1, 0 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx, columnIdx + 1 }
		startRowIdx = columnLength - 1
		startColumnIdx = 0
	case west:
		axisDirectionWidth = rowLenght
		tiltDirectionIdx = func(_, columnIdx int) int { return columnIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return 0, columnIdx + 1 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx + 1, columnIdx }
		startRowIdx = 0
		startColumnIdx = 0
	case east:
		axisDirectionWidth = rowLenght
		tiltDirectionIdx = func(_, columnIdx int) int { return columnIdx }
		updateTiltDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return 0, columnIdx - 1 }
		updateAxisDirectionIdx = func(rowIdx, columnIdx int) (int, int) { return rowIdx + 1, columnIdx }
		startRowIdx = 0
		startColumnIdx = rowLenght - 1
	default:
		panic(fmt.Sprintf("unknown dir: %d", dir))
	}

	var res int
	rowIdx, columnIdx := startRowIdx, startColumnIdx
	for {
		for {
			tile := d[rowIdx][columnIdx]
			switch tile {
			case roundedRock:
				res += (axisDirectionWidth - tiltDirectionIdx(rowIdx, columnIdx))
			}

			rowIdx, columnIdx = updateAxisDirectionIdx(rowIdx, columnIdx)
			if rowIdx >= columnLength {
				break
			}
			if columnIdx >= rowLenght {
				break
			}
		}
		rowIdx, columnIdx = updateTiltDirectionIdx(rowIdx, columnIdx)
		if rowIdx >= columnLength {
			break
		}
		if columnIdx >= rowLenght {
			break
		}
	}

	return res
}

func (d Dish) String() string {
	var sb strings.Builder
	for i := 0; i < len(d); i++ {
		for j := 0; j < len(d[i]); j++ {
			sb.WriteRune(rune(d[i][j]))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type Pos struct {
	rowIdx, columnIdx int
}
