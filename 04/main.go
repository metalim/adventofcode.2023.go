package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var reCard = regexp.MustCompile(`^Card\s+(\d+):(.*)\|(.*)$`)

func part1(lines []string) {
	timeStart := time.Now()
	var sum int
	for _, line := range lines {
		if line == "" {
			continue
		}

		m := reCard.FindStringSubmatch(line)
		if m == nil {
			panic("invalid line: " + line)
		}

		winning := map[string]bool{}
		for _, win := range strings.Fields(m[2]) {
			winning[win] = true
		}

		var points int
		for _, card := range strings.Fields(m[3]) {
			if winning[card] {
				if points == 0 {
					points = 1
				} else {
					points *= 2
				}
			}
		}

		sum += points
	}

	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	var total int
	copies := map[int]int{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		m := reCard.FindStringSubmatch(line)
		if m == nil {
			panic("invalid line: " + line)
		}

		card, err := strconv.Atoi(m[1])
		catch(err)
		numCards := 1 + copies[card]
		total += numCards

		winning := map[string]bool{}
		for _, win := range strings.Fields(m[2]) {
			winning[win] = true
		}

		var points int
		for _, card := range strings.Fields(m[3]) {
			if winning[card] {
				points++
			}
		}

		for n := card + 1; points > 0; {
			copies[n] += numCards
			n++
			points--
		}
	}
	fmt.Println("Part 2:", total, "\tin", time.Since(timeStart))
}
