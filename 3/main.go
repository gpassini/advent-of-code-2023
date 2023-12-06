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
	var res int
	var lineIdx int

	var previousLineNumbers []*NumberBuilder
	var previousLineGears []*Gear

	r := bufio.NewScanner(strings.NewReader(input))
	for r.Scan() {
		var lineNumbers []*NumberBuilder
		var lineGears []*Gear

		var currNumber *NumberBuilder
		for i, c := range r.Text() {
			if n, ok := toInt(c); ok {
				if currNumber == nil {
					currNumber = NewNumberBuilder(lineIdx)
				}
				currNumber.Consume(n, i)
			} else {
				// check if we had a number in progress
				if currNumber != nil {
					lineNumbers = append(lineNumbers, currNumber)
					currNumber = nil
				}

				if c == '*' {
					gear := NewGear(lineIdx, i)
					lineGears = append(lineGears, gear)
				}
			}
		}
		// check if we had a number in progress
		if currNumber != nil {
			lineNumbers = append(lineNumbers, currNumber)
		}

		// check if we have gears close to numbers (in this line and against the previous)
		for _, number := range lineNumbers {
			fmt.Println("Number:", number)
			for _, gear := range lineGears {
				if IsConnected(number, gear) {
					gear.Consume(number.value)
				}
			}
			for _, gear := range previousLineGears {
				if IsConnected(number, gear) {
					gear.Consume(number.value)
				}
			}
		}
		for _, gear := range lineGears {
			// only look at the previous line because we've already matched all gears with numbers from this line just above
			for _, number := range previousLineNumbers {
				if IsConnected(number, gear) {
					gear.Consume(number.value)
				}
			}
		}

		// check what gears from the previous line are valid are add their results before we move on
		for _, gear := range previousLineGears {
			fmt.Println("Previous line gear:", gear)
			if gear.IsValid() {
				res += gear.Ratio()
			}
		}

		previousLineNumbers = lineNumbers
		previousLineGears = lineGears
		lineIdx++
	}

	// we must check the very last line
	for _, gear := range previousLineGears {
		fmt.Println("Previous line gear:", gear)
		if gear.IsValid() {
			res += gear.Ratio()
		}
	}

	println(res)
}

type NumberBuilder struct {
	value          int
	lineIdx        int
	startColumnIdx int
	endColumnIdx   int
}

func NewNumberBuilder(lineIdx int) *NumberBuilder {
	return &NumberBuilder{
		lineIdx:        lineIdx,
		startColumnIdx: -1,
	}
}

func (builder *NumberBuilder) Consume(n int, columnIdx int) {
	builder.value = builder.value*10 + n
	if builder.startColumnIdx == -1 {
		builder.startColumnIdx = int(columnIdx)
	}
	builder.endColumnIdx = columnIdx
}

type Gear struct {
	numbers   []int
	lineIdx   int
	columnIdx int
}

func NewGear(lineIdx, columnIdx int) *Gear {
	return &Gear{
		lineIdx:   lineIdx,
		columnIdx: columnIdx,
	}
}

func (g *Gear) Consume(n int) {
	g.numbers = append(g.numbers, n)
}

func (g Gear) IsValid() bool {
	return len(g.numbers) == 2
}

func (g Gear) Ratio() int {
	if !g.IsValid() {
		return 0
	}
	return g.numbers[0] * g.numbers[1]
}

func IsConnected(n *NumberBuilder, g *Gear) bool {
	return n.startColumnIdx-1 <= g.columnIdx && n.endColumnIdx+1 >= g.columnIdx
}

func toInt(c rune) (int, bool) {
	if c >= '0' && c <= '9' {
		return int(c - '0'), true
	}
	return 0, false
}
