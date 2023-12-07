package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

var (
	//go:embed input.txt
	input string
)

func main() {
	var res int
	r := bufio.NewScanner(strings.NewReader(input))
	var cardCopies []int
	for lineIdx := 0; r.Scan(); lineIdx++ {
		line := r.Text()

		// each card has 1 copy by default
		if len(cardCopies) < lineIdx+1 {
			cardCopies = append(cardCopies, 1)
		} else {
			cardCopies[lineIdx]++
		}

		// card ID
		cardId, lineTail := readCardId(line)
		card := Card{id: cardId}
		var number string
		// winning numbers
		more := true
		for more {
			number, lineTail, more = readNumber(lineTail)
			card.AddWinningNumber(number)
		}
		// numbers
		more = true
		for more {
			number, lineTail, more = readNumber(lineTail)
			card.AddNumber(number)
		}
		// count winning numbers
		score := card.CountWinningNumbers()
		for i := lineIdx + 1; i <= lineIdx+score; i++ {
			if len(cardCopies) < i+1 {
				cardCopies = append(cardCopies, cardCopies[lineIdx])
			} else {
				cardCopies[i] += cardCopies[lineIdx]
			}
		}
		fmt.Printf("%s has %d winning numbers and %d copies\n", card, score, cardCopies[lineIdx])
		res += cardCopies[lineIdx]
	}
	println(res)
}

type Card struct {
	id             string
	winningNumbers []string
	numbers        []string
}

func (c *Card) AddWinningNumber(n string) {
	if n == "" {
		return
	}
	c.winningNumbers = append(c.winningNumbers, n)
}

func (c *Card) AddNumber(n string) {
	if n == "" {
		return
	}
	c.numbers = append(c.numbers, n)
}

func (c Card) CalculateScore() int {
	var score int
	slices.Sort(c.winningNumbers)
	for _, n := range c.numbers {
		if _, ok := slices.BinarySearch(c.winningNumbers, n); ok {
			if score == 0 {
				score = 1
			} else {
				score = score * 2
			}
		}
	}
	return score
}

func (c Card) CountWinningNumbers() int {
	var count int
	slices.Sort(c.winningNumbers)
	for _, n := range c.numbers {
		if _, ok := slices.BinarySearch(c.winningNumbers, n); ok {
			count++
		}
	}
	return count
}

func readCardId(line string) (id string, rest string) {
	const prefix = "Card "
	line = line[len(prefix):]
	for i, c := range line {
		if c == ' ' {
			continue
		} else if c >= '0' && c <= '9' {
			id = id + string(c)
		} else {
			if id == "" {
				panic(fmt.Sprintf("couldn't find card id in line: %s", line))
			}
			// do not return the ':' in the rest
			return id, line[i+1:]
		}
	}
	panic(fmt.Sprintf("empty line? %s", line))
}

func readNumber(line string) (number string, rest string, ok bool) {
	for i, c := range line {
		if c == ' ' {
			if number != "" {
				return number, line[i:], true
			}
			continue
		} else if c == '|' {
			// do not return the '|' in the rest
			return number, line[i+1:], false
		} else if c >= '0' && c <= '9' {
			number = number + string(c)
		} else {
			panic(fmt.Sprintf("unhandled char: %c", c))
		}
	}
	return number, "", false
}
