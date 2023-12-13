package main

import (
	"fmt"
	"os"
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

func sumReflection(lines []string) int {
	var sum int
	for i := range lines {
		up := i
		down := i + 1
		reflection := true
		var rows int
		for up >= 0 && down < len(lines) {
			if lines[up] != lines[down] {
				reflection = false
				break
			}
			rows++
			up--
			down++
		}
		if reflection && rows > 0 {
			sum += (i + 1) * 100
		}
	}

	for i := range lines[0] {
		left := i
		right := i + 1
		reflection := true
		var cols int
		for left >= 0 && right < len(lines[0]) {
			for _, line := range lines {
				if line[left] != line[right] {
					reflection = false
					break
				}
			}
			if !reflection {
				break
			}
			cols++
			left--
			right++
		}
		if reflection && cols > 0 {
			sum += i + 1
		}
	}
	return sum
}

func part1(lines []string) {
	timeStart := time.Now()
	var start, sum int
	for i, line := range lines {
		if line == "" {
			sum += sumReflection(lines[start:i])
			start = i + 1
		}
	}
	sum += sumReflection(lines[start:])

	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(timeStart))
}
