package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	one   = "one"
	two   = "two"
	three = "three"
	four  = "four"
	five  = "five"
	six   = "six"
	seven = "seven"
	eight = "eight"
	nine  = "nine"
)

var (
	//go:embed input.txt
	input string

	allNumbersStrs = []string{
		one,
		two,
		three,
		four,
		five,
		six,
		seven,
		eight,
		nine,
	}
)

func main() {
	var res int
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		var firstC rune
		var firstCPos int
		var lastC rune
		var lastCPos int
		var states []writtenNumberPos
		line := scanner.Text()
		for pos, c := range line {
			if c >= '0' && c <= '9' {
				if firstC == 0 {
					firstC = c
					firstCPos = pos
				} else {
					lastC = c
					lastCPos = pos
				}
			} else {
				tempStates := canBeNumber(c, pos, states...)
				states = []writtenNumberPos{}
				for _, s := range tempStates {
					if s.nextPos >= len(s.numberStr) {
						// complete number
						if firstC == 0 || firstCPos > s.linePos {
							firstC = rune(strconv.Itoa(s.number)[0])
							firstCPos = s.linePos
						} else if lastCPos < s.linePos {
							lastC = rune(strconv.Itoa(s.number)[0])
							lastCPos = s.linePos
						}
					} else {
						states = append(states, s)
					}
				}
				// fmt.Println("states at the end:", states)
			}
		}

		first, _ := strconv.Atoi(fmt.Sprintf("%c", firstC))
		if lastC == 0 {
			lastC = firstC
		}
		last, _ := strconv.Atoi(fmt.Sprintf("%c", lastC))
		val := first*10 + last
		fmt.Println(val)
		res += val
	}
	fmt.Println(res)
}

func canBeNumber(c rune, linePos int, states ...writtenNumberPos) []writtenNumberPos {
	// fmt.Println(fmt.Sprintf("%c", c), linePos, states)

	var nextStates []writtenNumberPos

	// check current states
	for _, s := range states {
		if c == rune(s.numberStr[s.nextPos]) {
			s.nextPos++
			nextStates = append(nextStates, s)
		}
	}

	// check if we can create new ones
	for n, numberStr := range allNumbersStrs {
		if c == rune(numberStr[0]) {
			nextStates = append(nextStates, writtenNumberPos{
				numberStr: numberStr,
				number:    n + 1,
				linePos:   linePos,
				nextPos:   1,
			})
		}
	}

	// fmt.Println(nextStates)

	return nextStates
}

type writtenNumberPos struct {
	numberStr string
	number    int
	linePos   int
	nextPos   int
}
