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

	// read seeds
	var available []int
	for r.Scan() {
		line := r.Text()
		indexOfSeeds := strings.Index(line, ":")
		if indexOfSeeds == -1 {
			panic(fmt.Sprintf("':' not found in line: %s", line))
		}
		rest := line[indexOfSeeds+1:]
		available = readAllNumbers(rest)
		break
	}
	fmt.Println("Available seeds:", available)

	for {
		encodedMap, ok := parseMap(r)
		if !ok {
			break
		}
		fmt.Println("Encoded map:", encodedMap)
		for i, s := range available {
			available[i] = encodedMap.FindDest(s)
		}
		fmt.Println("Steps:", available)
	}

	// find smallest value from last step
	minVal := available[0]
	for i := 1; i < len(available); i++ {
		if available[i] < minVal {
			minVal = available[i]
		}
	}

	println(minVal)
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

type EncodedMap struct {
	destStarts   []int
	sourceStarts []int
	ranges       []int
}

func (m *EncodedMap) Add(destStart, sourceStart, mapRange int) {
	m.destStarts = append(m.destStarts, destStart)
	m.sourceStarts = append(m.sourceStarts, sourceStart)
	m.ranges = append(m.ranges, mapRange)
}

func (m EncodedMap) FindDest(source int) int {
	for i := 0; i < len(m.sourceStarts); i++ {
		sourceStart := m.sourceStarts[i]
		destStart := m.destStarts[i]
		applicableRange := m.ranges[i]
		if source >= sourceStart && source <= sourceStart+applicableRange {
			diff := source - sourceStart
			return destStart + diff
		}
	}
	// unmmaped means same number
	return source
}

func parseMap(r *bufio.Scanner) (EncodedMap, bool) {
	m := EncodedMap{}
	for r.Scan() {
		// skip until after the header
		line := r.Text()
		if !strings.HasSuffix(line, ":") {
			continue
		}
		break
	}
	for r.Scan() {
		line := r.Text()
		if line == "" {
			break
		}
		mapCodes := readAllNumbers(line)
		if len(mapCodes) != 3 {
			panic(fmt.Sprintf("expected 3 numbers from line (got %d): %s", len(mapCodes), line))
		}
		m.Add(mapCodes[0], mapCodes[1], mapCodes[2])
	}
	if len(m.destStarts) == 0 {
		return m, false
	}
	return m, true
}
