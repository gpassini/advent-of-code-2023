package main

import (
	"bufio"
	_ "embed"
	"strings"
)

var (
	//go:embed input.txt
	input string
)

func main() {
	r := bufio.NewScanner(strings.NewReader(input))
	r.Scan()
	initSeq := r.Text()
	var res int
	steps := strings.Split(initSeq, ",")
	for _, step := range steps {
		res += hash(step)
	}
	println(res)
}

func hash(s string) int {
	const (
		mul = 17
		mod = 256
	)

	var res int
	for _, c := range s {
		res += int(c)
		res *= mul
		res %= mod
	}
	return res
}
