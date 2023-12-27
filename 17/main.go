package main

import (
	"bufio"
	_ "embed"
	"strings"
)

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	var cityMap CityMap
	for r.Scan() {
		var row []int
		for _, c := range r.Text() {
			row = append(row, int(c-'0'))
		}
		cityMap = append(cityMap, row)
	}

	var res int
	cruciblesToVisit := []CrucibleToVisit{
		{
			rowIdx:           0,
			columnIdx:        1,
			incomingDir:      right,
			dirStreak:        1,
			previousHeatLoss: 0,
		},
		{
			rowIdx:           1,
			columnIdx:        0,
			incomingDir:      down,
			dirStreak:        1,
			previousHeatLoss: 0,
		},
	}
	visitedCruciblesToMinHeatLoss := make(map[VisitedCrucible]int)
	for len(cruciblesToVisit) > 0 {
		newCruciblesToVisit, finalHeatLoss := cityMap.Visit(cruciblesToVisit[0], visitedCruciblesToMinHeatLoss)
		if finalHeatLoss > 0 {
			if res == 0 || finalHeatLoss < res {
				res = finalHeatLoss
			}
		}
		cruciblesToVisit = append(cruciblesToVisit[1:], newCruciblesToVisit...)
	}
	println(res)
}

type CityMap [][]int

func (m CityMap) Visit(ctv CrucibleToVisit, visited map[VisitedCrucible]int) ([]CrucibleToVisit, int) {
	rowIdx, columnIdx := ctv.rowIdx, ctv.columnIdx
	if rowIdx < 0 || rowIdx >= len(m) || columnIdx < 0 || columnIdx >= len(m[0]) {
		return nil, 0
	}

	incomingDir := ctv.incomingDir
	dirStreak := ctv.dirStreak
	heatLoss := ctv.previousHeatLoss + m[rowIdx][columnIdx]
	currCrucible := VisitedCrucible{
		rowIdx:      rowIdx,
		columnIdx:   columnIdx,
		incomingDir: incomingDir,
		streak:      dirStreak,
	}
	if minHeatLoss, ok := visited[currCrucible]; ok && heatLoss >= minHeatLoss {
		return nil, 0
	} else {
		visited[currCrucible] = heatLoss
	}

	if rowIdx == len(m)-1 && columnIdx == len(m[0])-1 {
		// arrived at destination, no need for further visits
		return nil, heatLoss
	}

	possibleDirections := incomingDir.PossibleDirections(dirStreak)
	var newCruciblesToVisit []CrucibleToVisit
	for _, dir := range possibleDirections {
		newRowIdx, newColumnIdx := dir.dir.Move(rowIdx, columnIdx)
		newCruciblesToVisit = append(newCruciblesToVisit, CrucibleToVisit{
			rowIdx:           newRowIdx,
			columnIdx:        newColumnIdx,
			incomingDir:      dir.dir,
			dirStreak:        dir.streak,
			previousHeatLoss: heatLoss,
		})
	}
	return newCruciblesToVisit, 0
}

type Coordinate struct {
	rowIdx, columnIdx int
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
	}
	return rowIdx, columnIdx
}

func (d Direction) RotateRight() Direction {
	switch d {
	case up:
		return right
	case down:
		return left
	case left:
		return up
	case right:
		return down
	}
	panic("invalid direction")
}

func (d Direction) RotateLeft() Direction {
	switch d {
	case up:
		return left
	case down:
		return right
	case left:
		return down
	case right:
		return up
	}
	panic("invalid direction")
}

func (d Direction) PossibleDirections(streak int) []PossibleDirection {
	dirs := []PossibleDirection{
		{
			dir:    d.RotateRight(),
			streak: 1,
		},
		{
			dir:    d.RotateLeft(),
			streak: 1,
		},
	}
	if streak < 3 {
		dirs = append(dirs, PossibleDirection{
			dir:    d,
			streak: streak + 1,
		})
	}
	return dirs
}

type PossibleDirection struct {
	dir    Direction
	streak int
}

type VisitedCrucible struct {
	rowIdx, columnIdx int
	incomingDir       Direction
	streak            int
}

type CrucibleToVisit struct {
	rowIdx, columnIdx int
	incomingDir       Direction
	dirStreak         int
	previousHeatLoss  int
}
