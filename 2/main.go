package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	red   = "red"
	green = "green"
	blue  = "blue"

	redLimit   = 12
	greenLimit = 13
	blueLimit  = 14
)

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	res := 0
	for r.Scan() {
		line := r.Text()
		game := parseLine(line)
		// valid := true
		// for _, s := range game.sets {
		// 	if s[red] > redLimit || s[blue] > blueLimit || s[green] > greenLimit {
		// 		valid = false
		// 		break
		// 	}
		// }
		// if valid {
		// fmt.Println(game)
		res += game.mins[red] * game.mins[green] * game.mins[blue]
		// }
	}
	println(res)
}

type game struct {
	id   int
	sets []set
	mins set
}

// color to count
type set map[string]int

func parseLine(l string) game {
	l = strings.TrimPrefix(l, "Game ")
	id, lastIdIdx := readNextNumber(l)
	game := game{
		id:   id,
		mins: map[string]int{red: 0, green: 0, blue: 0},
	}
	l = l[lastIdIdx+1:]
	l = strings.TrimPrefix(l, ":")
	setsStrs := strings.Split(l, ";")
	for _, setStr := range setsStrs {
		set := make(set)
		revelations := strings.Split(setStr, ",")
		for _, revelation := range revelations {
			revelation = strings.TrimPrefix(revelation, " ")
			count, lastIdx := readNextNumber(revelation)
			revelation = revelation[lastIdx+1:]
			color, _ := readNextColor(revelation)
			set[color] = count
			if count > game.mins[color] {
				game.mins[color] = count
			}
		}
		game.sets = append(game.sets, set)
	}
	return game
}

func readNextNumber(s string) (res int, lastIdx int) {
	for i, c := range s {
		lastIdx = i
		if c >= '0' && c <= '9' {
			n, err := strconv.Atoi(fmt.Sprintf("%c", c))
			if err != nil {
				panic(fmt.Sprintf("failed to parse next number (s: %s | c: %c)", s, c))
			}
			res = res*10 + n
		} else {
			break
		}
	}
	return res, lastIdx
}

func readNextColor(s string) (color string, lastIdx int) {
	for i, c := range s {
		switch c {
		case 'r':
			return red, i + len(red) - 1
		case 'g':
			return green, i + len(green) - 1
		case 'b':
			return blue, i + len(blue) - 1
		}
	}
	panic(fmt.Sprintf("no color found in: %s", s))
}
