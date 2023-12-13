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

func transpose(lines []string) []string {
	out := make([]string, len(lines[0]))
	for i := range out {
		for _, line := range lines {
			out[i] += string(line[i])
		}
	}
	return out
}

func part1(lines []string) {
	timeStart := time.Now()
	var start, sum int
	for i, line := range lines {
		if line == "" {
			sum += sumReflections(lines[start:i], 0)
			start = i + 1
		}
	}
	sum += sumReflections(lines[start:], 0)

	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func getDiff(a, b string) int {
	var diff int
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff
}

func sumRowsReflections(lines []string, diff int) int {
	var sum int
	for i := range lines {
		up := i
		down := i + 1
		var rows, smudges int
		for up >= 0 && down < len(lines) {
			smudges += getDiff(lines[up], lines[down])
			if smudges > diff {
				break
			}
			rows++
			up--
			down++
		}
		if rows > 0 && smudges == diff {
			sum += i + 1
		}
	}
	return sum
}

func sumReflections(lines []string, diff int) int {
	return sumRowsReflections(lines, diff)*100 + sumRowsReflections(transpose(lines), diff)
}

func part2(lines []string) {
	timeStart := time.Now()
	var start, sum int
	_ = start
	for i, line := range lines {
		if line == "" {
			sum += sumReflections(lines[start:i], 1)
			start = i + 1
		}
	}
	sum += sumReflections(lines[start:], 1)

	fmt.Println("Part 2:", sum, "\tin", time.Since(timeStart))
}
