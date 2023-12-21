package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

const (
	none Orientation = iota
	vertical
	horizontal
)

type Orientation uint

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	var patterns []Pattern
	currPattern := Pattern{
		lineToIdxs:   make(map[string][]int),
		columnToIdxs: make(map[string][]int),
	}
	for r.Scan() {
		line := r.Text()
		if line == "" {
			for j, c := range currPattern.columns {
				currPattern.columnToIdxs[c] = append(currPattern.columnToIdxs[c], j)
			}
			patterns = append(patterns, currPattern)
			currPattern = Pattern{
				lineToIdxs:   make(map[string][]int),
				columnToIdxs: make(map[string][]int),
			}
			continue
		}

		currPattern.lines = append(currPattern.lines, line)
		currPattern.lineToIdxs[line] = append(currPattern.lineToIdxs[line], len(currPattern.lines)-1)

		if len(currPattern.columns) == 0 {
			currPattern.columns = make([]string, len(line))
		}
		for i, c := range line {
			currPattern.columns[i] += string(c)
		}
	}
	for j, c := range currPattern.columns {
		currPattern.columnToIdxs[c] = append(currPattern.columnToIdxs[c], j)
	}
	patterns = append(patterns, currPattern)

	var res int
	for i, p := range patterns {
		fmt.Println("Pattern", i, p)
		orientation, idx := p.SearchReflection()
		switch orientation {
		case vertical:
			res += idx + 1
		case horizontal:
			res += 100 * (idx + 1)
		}
	}
	println(res)
}

type Pattern struct {
	lines   []string
	columns []string

	lineToIdxs   map[string][]int
	columnToIdxs map[string][]int
}

func (p Pattern) SearchReflection() (Orientation, int) {
	for i := 0; i < len(p.lines)-1; i++ {
		found := true
		for iBack, iFor := i, i+1; iBack >= 0 && iFor < len(p.lines); iBack, iFor = iBack-1, iFor+1 {
			if p.lines[iBack] != p.lines[iFor] {
				found = false
				break
			}
		}
		if found {
			return horizontal, i
		}
	}

	for i := 0; i < len(p.columns)-1; i++ {
		found := true
		for iBack, iFor := i, i+1; iBack >= 0 && iFor < len(p.columns); iBack, iFor = iBack-1, iFor+1 {
			if p.columns[iBack] != p.columns[iFor] {
				found = false
				break
			}
		}
		if found {
			return vertical, i
		}
	}

	return none, 0
}

func (p Pattern) String() string {
	var sb strings.Builder
	sb.WriteString("Lines:\n")
	for _, l := range p.lines {
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	sb.WriteString("Columns:\n")
	for _, c := range p.columns {
		sb.WriteString(c)
		sb.WriteString("\n")
	}
	return sb.String()
}
