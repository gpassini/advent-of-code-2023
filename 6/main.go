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
	r.Scan()
	times := readAllNumbers(r.Text())
	fmt.Println("Times:", times)
	r.Scan()
	distances := readAllNumbers(r.Text())
	fmt.Println("Distances:", distances)
	for i, time := range times {
		n := calculateOptions(time, distances[i])
		if res == 0 {
			res = n
		} else {
			res *= n
		}
	}
	println(res)
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
