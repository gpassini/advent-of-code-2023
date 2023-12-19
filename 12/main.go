package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	operational SpringStatus = '.'
	damaged     SpringStatus = '#'
	unknown     SpringStatus = '?'
)

type SpringStatus rune

func (s SpringStatus) String() string {
	return string(rune(s))
}

var (
	//go:embed input.txt
	input string
)

func main() {
	var res int
	r := bufio.NewScanner(strings.NewReader(input))
	for r.Scan() {
		record := parseRecord(r.Text())
		fmt.Println(record)
		pos := record.calculatePossibilities(0)
		fmt.Println(pos)
		res += pos
	}
	println(res)

	// record := parseRecord("?? 1,1")
	// fmt.Println(record)
	// pos := record.calculatePossibilities(0)
	// fmt.Println(pos)
}

func parseRecord(s string) Record {
	record := Record{}

	var lastIdx int
	for i, c := range s {
		lastIdx = i
		if c == ' ' {
			break
		}
		record.springs = append(record.springs, SpringStatus(c))
	}

	var currNumber int
	for _, c := range s[lastIdx+1:] {
		if c >= '0' && c <= '9' {
			currNumber = currNumber*10 + int(c-'0')
		} else {
			record.damagedSpringGroups = append(record.damagedSpringGroups, currNumber)
			currNumber = 0
		}
	}
	record.damagedSpringGroups = append(record.damagedSpringGroups, currNumber)

	return record
}

type Record struct {
	springs             []SpringStatus
	damagedSpringGroups []int
}

func (r Record) calculatePossibilities(damagedStreak int) int {
	if len(r.springs) == 0 {
		if len(r.damagedSpringGroups) == 0 ||
			(len(r.damagedSpringGroups) == 1 && r.damagedSpringGroups[0] == damagedStreak) {
			return 1
		} else {
			return 0
		}
	}

	spring := r.springs[0]
	switch spring {
	case operational:
		if damagedStreak > 0 {
			if damagedStreak != r.damagedSpringGroups[0] {
				return 0
			} else {
				r.damagedSpringGroups = r.damagedSpringGroups[1:]
			}
		}
		r.springs = r.springs[1:]
		return r.calculatePossibilities(0)
	case damaged:
		r.springs = r.springs[1:]
		newDamageStreak := damagedStreak + 1
		if len(r.damagedSpringGroups) == 0 || r.damagedSpringGroups[0] < newDamageStreak {
			return 0
		}
		return r.calculatePossibilities(newDamageStreak)
	case unknown:
		rWithOperational := Record{
			springs:             append([]SpringStatus{operational}, r.springs[1:]...),
			damagedSpringGroups: r.damagedSpringGroups,
		}
		rWithDamaged := Record{
			springs:             append([]SpringStatus{damaged}, r.springs[1:]...),
			damagedSpringGroups: r.damagedSpringGroups,
		}
		return rWithOperational.calculatePossibilities(damagedStreak) +
			rWithDamaged.calculatePossibilities(damagedStreak)
	default:
		panic(fmt.Sprintf("unknown spring status: %s", spring))
	}
}

func (r Record) String() string {
	var sb strings.Builder
	for _, s := range r.springs {
		sb.WriteRune(rune(s))
	}
	sb.WriteRune(' ')
	for _, g := range r.damagedSpringGroups {
		sb.WriteString(strconv.Itoa(g))
		sb.WriteRune(',')
	}
	return sb.String()
}
