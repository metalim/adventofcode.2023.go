package main

import (
	"fmt"
	"os"
	"strings"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(os.Args[1])
	catch(err)
	lines := strings.Split(string(bs), "\n")

	part1(lines)
	part2(lines)
}

func part1(lines []string) {
	var sum int
	for _, line := range lines {
		var digits, first, last int
		for _, c := range line {
			if '0' <= c && c <= '9' {
				digit := int(c - '0')
				if digits == 0 {
					first = digit
				}
				last = digit
				digits++
			}
		}
		sum += first*10 + last
	}
	fmt.Println(sum)
}

func part2(lines []string) {
	var sum int
	for _, line := range lines {
		var digits, first, last int
		for i := range line {
			digit, ok := getDigit(line[i:])
			if !ok {
				continue
			}
			if digits == 0 {
				first = digit
			}
			last = digit
			digits++
		}
		sum += first*10 + last
	}
	fmt.Println(sum)
}

func getDigit(s string) (int, bool) {
	if '0' <= s[0] && s[0] <= '9' {
		return int(s[0] - '0'), true
	}

	// no "zero"
	names := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	for i, name := range names {
		if strings.HasPrefix(s, name) {
			return i + 1, true
		}
	}
	return 0, false
}
