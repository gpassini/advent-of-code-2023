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
	r := bufio.NewScanner(strings.NewReader(input))
	r.Scan()
	time := readOneNumber(r.Text())
	fmt.Println("Time:", time)
	r.Scan()
	distance := readOneNumber(r.Text())
	fmt.Println("Distance:", distance)
	println(calculateOptions(time, distance))
}

func calculateOptions(time, distance int) int {
	var count int
	for i := 1; i <= time-i; i++ {
		actual := i * (time - i)
		if actual > distance {
			count++
			if i != time-i {
				count++ // the inverse is also an option
			}
		}
	}
	return count
}

func readOneNumber(s string) int {
	var currNumber int
	for _, c := range s {
		if n, ok := charIsInt(c); ok {
			currNumber = currNumber*10 + n
		}
	}
	return currNumber
}

func readAllNumbers(s string) []int {
	var numbers []int
	var currNumber int
	var readingNumber bool
	for _, c := range s {
		if n, ok := charIsInt(c); ok {
			currNumber = currNumber*10 + n
			readingNumber = true
		} else {
			if readingNumber {
				numbers = append(numbers, currNumber)
				currNumber = 0
				readingNumber = false
			}
		}
	}
	if readingNumber {
		numbers = append(numbers, currNumber)
	}
	return numbers
}

func charIsInt(c rune) (int, bool) {
	if c >= '0' && c <= '9' {
		return int(c - '0'), true
	} else {
		return 0, false
	}
}
