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

func (o Orientation) String() string {
	switch o {
	case none:
		return "none"
	case vertical:
		return "vertical"
	case horizontal:
		return "horizontal"
	default:
		panic(fmt.Sprintf("unkown orientation: %d", uint(o)))
	}
}

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	var patterns []Pattern
	currPattern := Pattern{}
	for r.Scan() {
		line := r.Text()
		if line == "" {
			patterns = append(patterns, currPattern)
			currPattern = Pattern{}
			continue
		}

		currPattern.lines = append(currPattern.lines, line)

		if len(currPattern.columns) == 0 {
			currPattern.columns = make([]string, len(line))
		}
		for i, c := range line {
			currPattern.columns[i] += string(c)
		}
	}
	patterns = append(patterns, currPattern)

	var res int
	for i, p := range patterns {
		fmt.Println("Pattern", i, p)
		oldOrientation, oldIdx := p.SearchReflection(none, -1)
		fmt.Println("Old reflection:", oldOrientation, oldIdx)
		orientation, idx := p.FixSmudge(oldOrientation, oldIdx)
		switch orientation {
		case vertical:
			res += idx + 1
		case horizontal:
			res += 100 * (idx + 1)
		case none:
			panic("smudge not found")
		}
	}
	println(res)
}

type Pattern struct {
	lines   []string
	columns []string
}

func (p Pattern) FixSmudge(oldOrientation Orientation, oldIdx int) (Orientation, int) {
	for i := 0; i < len(p.lines)-1; i++ {
		iLine := p.lines[i]
		for j := i + 1; j < len(p.lines); j++ {
			jLine := p.lines[j]
			if idx, ok := MaybeSmudge(iLine, jLine); ok {
				iLine := ChangeLine(iLine, idx)
				iPattern := Pattern{
					lines:   InsertAtIdx(p.lines, iLine, i),
					columns: p.columns,
				}
				if o, rIdx := iPattern.SearchReflection(oldOrientation, oldIdx); o != none {
					fmt.Printf("Smudge at line (%d,%d) with new reflection line at %d\n", i, idx, rIdx)
					return o, rIdx
				}

				jLine := ChangeLine(jLine, idx)
				jPattern := Pattern{
					lines:   InsertAtIdx(p.lines, jLine, j),
					columns: p.columns,
				}
				if o, rIdx := jPattern.SearchReflection(oldOrientation, oldIdx); o != none {
					fmt.Printf("Smudge at line (%d,%d) with new reflection line at %d\n", j, idx, rIdx)
					return o, rIdx
				}
			}
		}
	}

	for i := 0; i < len(p.columns)-1; i++ {
		iColumn := p.columns[i]
		for j := i + 1; j < len(p.columns); j++ {
			jColumn := p.columns[j]
			fmt.Println("Comparing columns", i, iColumn, j, jColumn)
			if idx, ok := MaybeSmudge(iColumn, jColumn); ok {
				fmt.Println("Maybe smudge at", idx)
				newIColumn := ChangeLine(iColumn, idx)
				iPattern := Pattern{
					lines:   p.lines,
					columns: InsertAtIdx(p.columns, newIColumn, i),
				}
				if o, rIdx := iPattern.SearchReflection(oldOrientation, oldIdx); o != none {
					fmt.Printf("Smudge at column (%d,%d) with new reflection line at %d\n", i, idx, rIdx)
					return o, rIdx
				}

				newJColumn := ChangeLine(jColumn, idx)
				jPattern := Pattern{
					lines:   p.lines,
					columns: InsertAtIdx(p.columns, newJColumn, j),
				}
				if o, rIdx := jPattern.SearchReflection(oldOrientation, oldIdx); o != none {
					fmt.Printf("Smudge at column (%d,%d) with new reflection line at %d\n", j, idx, rIdx)
					return o, rIdx
				}
			}
		}
	}

	return none, -1
}

func InsertAtIdx[T any](s []T, e T, idx int) []T {
	newS := make([]T, len(s))
	for i := 0; i < len(s); i++ {
		if i == idx {
			newS[i] = e
		} else {
			newS[i] = s[i]
		}
	}
	return newS
}

func ChangeLine(s string, i int) string {
	bytes := make([]byte, len([]byte(s)))
	copy(bytes, []byte(s))
	toChange := bytes[i]
	if toChange == byte('#') {
		toChange = byte('.')
	} else {
		toChange = '#'
	}
	bytes[i] = toChange
	return string(bytes)
}

func MaybeSmudge(s1, s2 string) (int, bool) {
	idx := -1
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			if idx >= 0 {
				return 0, false
			}
			idx = i
		}
	}
	if idx == -1 {
		// identical lines
		return 0, false
	}
	return idx, true
}

func (p Pattern) SearchReflection(oldOrientation Orientation, oldIdx int) (Orientation, int) {
	for i := 0; i < len(p.lines)-1; i++ {
		found := true
		for iBack, iFor := i, i+1; iBack >= 0 && iFor < len(p.lines); iBack, iFor = iBack-1, iFor+1 {
			if p.lines[iBack] != p.lines[iFor] {
				found = false
				break
			}
		}
		if found && (oldOrientation != horizontal || oldIdx != i) {
			return horizontal, i
		}
	}

	for i := 0; i < len(p.columns)-1; i++ {
		found := true
		for iBack, iFor := i, i+1; iBack >= 0 && iFor < len(p.columns); iBack, iFor = iBack-1, iFor+1 {
			if p.columns[iBack] != p.columns[iFor] {
				fmt.Println("Columns not equal", iBack, p.columns[iBack], iFor, p.columns[iFor])
				found = false
				break
			}
		}
		if found {
			if oldOrientation != vertical || oldIdx != i {
				return vertical, i
			} else {
				fmt.Println("Discarting same reflection line", vertical, i)
			}
		}
	}

	return none, -1
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
