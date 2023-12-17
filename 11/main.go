package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

const (
	galaxyChar = '#'
)

var (
	//go:embed input.txt
	input string
)

func main() {
	var columnHasGalaxy []bool
	var lineHasGalaxy []bool
	var space Space

	r := bufio.NewScanner(strings.NewReader(input))
	var lineIdx int
	for r.Scan() {
		line := r.Text()

		lineHasGalaxy = append(lineHasGalaxy, false)

		for columnIdx, c := range line {
			if len(columnHasGalaxy) <= columnIdx {
				columnHasGalaxy = append(columnHasGalaxy, false)
			}

			if c == galaxyChar {
				space.AddGalaxy(lineIdx, columnIdx)
				lineHasGalaxy[lineIdx] = true
				columnHasGalaxy[columnIdx] = true
			}
		}

		lineIdx++
	}

	var res int

	for i, iGalaxy := range space.galaxies {
		for j := i + 1; j < len(space.galaxies); j++ {
			jGalaxy := space.galaxies[j]

			lineDistance := calculteDistance(iGalaxy.x, jGalaxy.x, lineHasGalaxy)
			fmt.Println("Line distance:", lineDistance)
			columnDistance := calculteDistance(iGalaxy.y, jGalaxy.y, columnHasGalaxy)
			fmt.Println("Column distance", columnDistance)

			distance := lineDistance + columnDistance

			fmt.Printf("Distance between %v and %v: %d\n", iGalaxy, jGalaxy, distance)

			res += distance
		}
	}

	println(res)
}

func calculteDistance(start, end int, hasGalaxy []bool) int {
	min, max := start, end
	if min > max {
		min, max = max, min
	}
	distance := max - min
	for i := min; i < max; i++ {
		if !hasGalaxy[i] {
			distance += 999_999
		}
	}
	return distance
}

type Space struct {
	galaxies []Point
}

func (s *Space) AddGalaxy(x, y int) {
	s.galaxies = append(s.galaxies, Point{x, y})
}

type Point struct {
	x, y int
}
