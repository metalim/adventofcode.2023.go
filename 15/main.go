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

func HASH(s string) int {
	var h int
	for _, c := range s {
		h = (h + int(c)) * 17 % 256
	}
	return h
}

func part1(lines []string) {
	timeStart := time.Now()
	var sum int
	for _, line := range lines {
		for _, s := range strings.Split(line, ",") {
			sum += HASH(s)
		}
	}

	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(timeStart))
}
