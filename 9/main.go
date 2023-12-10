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
	r := bufio.NewScanner(strings.NewReader(input))
	for r.Scan() {
		h := parseLine(r.Text())
		fmt.Println("History:", h)
		e := h.Extrapolate()
		fmt.Println("Extrapolation:", e)
		res += e
	}
	println(res)
}

func parseLine(s string) History {
	var h History
	var n int
	rest := s
	for {
		if rest == "" {
			break
		}
		n, rest = readNextInt(rest)
		h = append(History{n}, h...)
	}
	return h
}

func readNextInt(s string) (int, string) {
	var n int
	var reading bool
	var negative bool
	var lastIdx int
	for i, c := range s {
		lastIdx = i
		if c >= '0' && c <= '9' {
			reading = true
			n = n*10 + int(c-'0')
		} else if c == '-' {
			if reading {
				panic("- while reading number")
			}
			negative = true
		} else if c == ' ' {
			if reading {
				break
			}
		}
	}
	if !reading {
		panic("no number found")
	}
	if negative {
		n = -n
	}
	return n, s[lastIdx+1:]
}

type History []int

func (h History) Extrapolate() int {
	var stacks []History
	currStack := h
	for {
		stacks = append(stacks, currStack)
		if currStack.allZeroes() {
			break
		}
		nextStack := make(History, len(currStack)-1)
		for i := 0; i < len(nextStack); i++ {
			nextStack[i] = currStack[i+1] - currStack[i]
		}
		currStack = nextStack
	}

	// start from the first non-zero history starting from last
	for i := len(stacks) - 2; i >= 0; i-- {
		currStack := stacks[i]
		deeperStack := stacks[i+1]
		newVal := currStack[len(currStack)-1] + deeperStack[len(deeperStack)-1]
		if i == 0 {
			return newVal
		}
		stacks[i] = append(currStack, newVal)
	}

	return h[len(h)-1]
}

func (h History) allZeroes() bool {
	for _, v := range h {
		if v != 0 {
			return false
		}
	}
	return true
}
