package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

const (
	highCard Combination = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

const (
	two Card = iota
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	jack
	queen
	king
	ace
)

type Combination int
type Card int

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	var hands []Hand
	for r.Scan() {
		hands = append(hands, parseLine(r.Text()))
	}
	slices.SortFunc(hands, func(a, b Hand) int {
		if diff := int(a.combination - b.combination); diff != 0 {
			return diff
		}
		for i := 0; i < 5; i++ {
			if diff := int(a.cards[i] - b.cards[i]); diff != 0 {
				return diff
			}
		}
		return 0
	})
	var res int
	for i, h := range hands {
		fmt.Println(h)
		res += h.bid * (i + 1)
	}
	println(res)
}

func parseLine(line string) Hand {
	var cards [5]rune
	var bid int
	for i, c := range line {
		if i < 5 {
			cards[i] = c
		} else if c >= '0' && c <= '9' {
			bid = bid*10 + int(c-'0')
		}
	}
	return NewHand(cards, bid)
}

type Hand struct {
	cards       [5]Card
	bid         int
	combination Combination
}

func NewHand(cards [5]rune, bid int) Hand {
	var parsedCards [5]Card
	for i, c := range cards {
		parsedCards[i] = charToCard(c)
	}
	return Hand{
		cards:       parsedCards,
		bid:         bid,
		combination: calculateCombination(parsedCards),
	}
}

func calculateCombination(cards [5]Card) Combination {
	cardToCount := make(map[Card]int, len(cards))
	for _, card := range cards {
		cardToCount[card] += 1
	}

	fmt.Printf("Cards: %v\nCounts: %v\n", cards, cardToCount)

	var highestCount int
	var secondHighestCount int
	for _, count := range cardToCount {
		if count > highestCount {
			highestCount, count = count, highestCount
		}
		if count > secondHighestCount {
			secondHighestCount = count
		}
	}

	fmt.Println(highestCount, secondHighestCount)

	switch highestCount {
	case 5:
		return fiveOfAKind
	case 4:
		return fourOfAKind
	case 3:
		if secondHighestCount == 2 {
			return fullHouse
		} else {
			return threeOfAKind
		}
	case 2:
		if secondHighestCount == 2 {
			return twoPair
		} else {
			return onePair
		}
	default:
		return highCard
	}
}

func charToCard(c rune) Card {
	if c >= '2' && c <= '9' {
		return Card(c - '0' - 2)
	} else if c == 'T' {
		return ten
	} else if c == 'J' {
		return jack
	} else if c == 'Q' {
		return queen
	} else if c == 'K' {
		return king
	} else if c == 'A' {
		return ace
	}
	panic(fmt.Sprintf("non card char: %c", c))
}
