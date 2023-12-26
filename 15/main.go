package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	mul = 17
	mod = 256
)

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	r.Scan()
	initSeq := r.Text()
	steps := strings.Split(initSeq, ",")
	boxes := make(Boxes, mod)
	for _, step := range steps {
		var label string
		var focalLength int
		var remove bool
		splits := strings.Split(step, "=")
		if len(splits) == 2 {
			label = splits[0]
			focalLength, _ = strconv.Atoi(splits[1])
		} else {
			// remove the trailing '-'
			label = step[:len(step)-1]
			remove = true
		}

		labelHash := hash(label)

		if remove {
			boxes[labelHash] = removeLens(boxes[labelHash], label)
		} else {
			boxes[labelHash] = addLens(boxes[labelHash], Lens{label: label, focalLength: focalLength})
		}
		fmt.Println(boxes)
		println(boxes.FocusingPower())
	}
}

type Boxes map[int][]Lens

func (b Boxes) FocusingPower() int {
	var res int
	for boxIdx := 0; boxIdx < mod; boxIdx++ {
		lenses, _ := b[boxIdx]
		for lensIdx, lens := range lenses {
			lensFocusingPower := (boxIdx + 1) * (lensIdx + 1) * lens.focalLength
			fmt.Println(lens.label, lensFocusingPower)
			res += lensFocusingPower
		}
	}
	return res
}

func (b Boxes) String() string {
	var sb strings.Builder
	for i := 0; i < mod; i++ {
		lenses, _ := b[i]
		if len(lenses) > 0 {
			sb.WriteString(fmt.Sprintf("Box %d: ", i))
			for _, lens := range lenses {
				sb.WriteString(fmt.Sprintf("[%s %d] ", lens.label, lens.focalLength))
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

type Lens struct {
	label       string
	focalLength int
}

func hash(s string) int {
	var res int
	for _, c := range s {
		res += int(c)
		res *= mul
		res %= mod
	}
	return res
}

func removeLens(s []Lens, label string) []Lens {
	for i, lens := range s {
		if lens.label == label {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func addLens(s []Lens, lensToAdd Lens) []Lens {
	for i, lens := range s {
		if lens.label == lensToAdd.label {
			s[i] = lensToAdd
			return s
		}
	}
	return append(s, lensToAdd)
}
