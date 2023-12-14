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

func part1(lines []string) {
	timeStart := time.Now()
	var mmap [][]rune
	for y0, line := range lines {
		mmap = append(mmap, make([]rune, len(line)))
		for x, c := range line {
			if c == '.' || c == '#' {
				mmap[y0][x] = c
				continue
			}
			y := y0
			for y > 0 && mmap[y-1][x] == '.' {
				y--
			}
			mmap[y0][x] = '.'
			mmap[y][x] = c
		}
	}

	plot(mmap)
	var sum int
	for y, row := range mmap {
		for _, c := range row {
			if c == 'O' {
				sum += len(mmap) - y
			}
		}
	}
	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func plot(mmap [][]rune) {
	for y, row := range mmap {
		fmt.Printf("%2d: %s\n", len(mmap)-y, string(row))
	}
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(timeStart))
}
