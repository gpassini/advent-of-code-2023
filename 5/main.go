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
	var availableRanges []Range
	for r.Scan() {
		line := r.Text()
		indexOfSeeds := strings.Index(line, ":")
		if indexOfSeeds == -1 {
			panic(fmt.Sprintf("':' not found in line: %s", line))
		}
		rest := line[indexOfSeeds+1:]
		seedRanges := readAllNumbers(rest)
		for i := 0; i < len(seedRanges); i += 2 {
			rangeStart := seedRanges[i]
			rangeLength := seedRanges[i+1]
			availableRanges = append(availableRanges, Range{
				start: rangeStart,
				end:   rangeStart + rangeLength - 1,
			})
		}
		break
	}
	fmt.Println("Available seed ranges:", availableRanges)

	for {
		encodedMap, ok := parseMap(r)
		if !ok {
			break
		}
		fmt.Println("Encoded map:", encodedMap)
		availableRanges = encodedMap.FindPossibleRanges(availableRanges)
		var count int
		for _, r := range availableRanges {
			count += r.Size()
		}
		fmt.Println(count, "steps:", availableRanges)
	}

	// find smallest value from last step
	minVal := availableRanges[0].start
	for i := 1; i < len(availableRanges); i++ {
		if availableRanges[i].start < minVal {
			minVal = availableRanges[i].start
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
	fragments []MapFragment
}

type MapFragment struct {
	destStart   int
	sourceStart int
	rangeLength int
}

func (f MapFragment) Transform(source int) int {
	if source < f.sourceStart && source > f.sourceStart+f.rangeLength-1 {
		panic(fmt.Sprintf("source %d is out of bounds of fragment: %v", source, f))
	}
	return f.destStart + source - f.sourceStart
}

func (f MapFragment) SplitRange(r Range) (*Range, []Range) {
	fmt.Println(f, r)

	fragStart := f.sourceStart
	fragEnd := fragStart + f.rangeLength - 1
	rangeStart := r.start
	rangeEnd := r.end

	if rangeStart < fragStart {
		if rangeEnd < fragStart {
			// the range is completely out of this fragment
			return nil, []Range{r}
		} else if rangeEnd <= fragEnd {
			// the range includes the beginning of this fragment
			return NewRangePointer(f.Transform(fragStart), f.Transform(rangeEnd)),
				[]Range{NewRange(rangeStart, fragStart-1)}
		} else { // rangeEnd > fragEnd
			// the range incudes the whole fragment
			return NewRangePointer(f.Transform(fragStart), f.Transform(fragEnd)),
				[]Range{
					NewRange(rangeStart, fragStart),
					NewRange(fragEnd+1, rangeEnd),
				}
		}
	} else if rangeStart <= fragEnd {
		// the range starts in the middle of this fragment
		if rangeEnd <= fragEnd {
			// the range is completely included in the fragment
			return NewRangePointer(f.Transform(rangeStart), f.Transform(rangeEnd)),
				[]Range{}
		} else { // rangeEnd > fragEnd
			// the range includes the end of this fragment
			return NewRangePointer(f.Transform(rangeStart), f.Transform(fragEnd)),
				[]Range{NewRange(fragEnd+1, rangeEnd)}
		}
	} else {
		// the range is completely out of this fragment
		return nil, []Range{r}
	}
}

func (m *EncodedMap) Add(destStart, sourceStart, mapRange int) {
	m.fragments = append(m.fragments, MapFragment{
		destStart:   destStart,
		sourceStart: sourceStart,
		rangeLength: mapRange,
	})
}

func (m EncodedMap) FindDest(source int) int {
	for _, frag := range m.fragments {
		sourceStart := frag.sourceStart
		destStart := frag.destStart
		applicableRange := frag.rangeLength
		if source >= sourceStart && source <= sourceStart+applicableRange {
			diff := source - sourceStart
			return destStart + diff
		}
	}
	// unmmaped means same number
	return source
}

func (m EncodedMap) FindPossibleRanges(availableRanges []Range) []Range {
	var possibleRanges []Range
	rangesToHandle := availableRanges
	for _, frag := range m.fragments {
		var unhandledRanges []Range
		for i := 0; i < len(rangesToHandle); i++ {
			currRange := rangesToHandle[i]
			head, tail := frag.SplitRange(currRange)
			if head != nil {
				possibleRanges = append(possibleRanges, *head)
			}
			if len(tail) > 0 {
				unhandledRanges = append(unhandledRanges, tail...)
			}
		}
		rangesToHandle = unhandledRanges
	}
	return append(possibleRanges, rangesToHandle...)
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
	if len(m.fragments) == 0 {
		return m, false
	}
	return m, true
}

type Range struct {
	start int
	end   int // inclusive
}

func NewRange(start, end int) Range {
	if start > end {
		panic(fmt.Sprintf("start: %d, end: %d", start, end))
	}
	return Range{
		start: start,
		end:   end,
	}
}

func NewRangePointer(start, end int) *Range {
	r := NewRange(start, end)
	return &r
}

func (r Range) Size() int {
	size := r.end - r.start + 1
	if size <= 0 {
		panic(fmt.Sprintf("%v", r))
	}
	return size
}
