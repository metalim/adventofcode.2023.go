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

func parseMap(lines []string) [][]rune {
	mmap := make([][]rune, len(lines))
	for i, line := range lines {
		mmap[i] = []rune(line)
	}
	return mmap
}

func part1(lines []string) {
	timeStart := time.Now()

	mmap := parseMap(lines)
	rollUp(mmap)

	// plot(mmap)
	sum := calcSum(mmap)
	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func calcSum(mmap [][]rune) int {
	var sum int
	for y, row := range mmap {
		for _, c := range row {
			if c == 'O' {
				sum += len(mmap) - y
			}
		}
	}
	return sum
}

func rollUp(mmap [][]rune) {
	for x := 0; x < len(mmap); x++ {
		for y0 := 0; y0 < len(mmap); y0++ {
			if mmap[y0][x] != 'O' {
				continue
			}
			y := y0
			for y > 0 && mmap[y-1][x] == '.' {
				y--
			}
			if y != y0 {
				mmap[y0][x] = '.'
				mmap[y][x] = 'O'
			}
		}
	}
}

func rollDown(mmap [][]rune) {
	for x := 0; x < len(mmap); x++ {
		for y0 := len(mmap) - 1; y0 >= 0; y0-- {
			if mmap[y0][x] != 'O' {
				continue
			}
			y := y0
			for y < len(mmap)-1 && mmap[y+1][x] == '.' {
				y++
			}
			if y != y0 {
				mmap[y0][x] = '.'
				mmap[y][x] = 'O'
			}
		}
	}
}

func rollLeft(mmap [][]rune) {
	for y := 0; y < len(mmap); y++ {
		for x0 := 0; x0 < len(mmap); x0++ {
			if mmap[y][x0] != 'O' {
				continue
			}
			x := x0
			for x > 0 && mmap[y][x-1] == '.' {
				x--
			}
			if x != x0 {
				mmap[y][x0] = '.'
				mmap[y][x] = 'O'
			}
		}
	}
}

func rollRight(mmap [][]rune) {
	for y := 0; y < len(mmap); y++ {
		for x0 := len(mmap) - 1; x0 >= 0; x0-- {
			if mmap[y][x0] != 'O' {
				continue
			}
			x := x0
			for x < len(mmap)-1 && mmap[y][x+1] == '.' {
				x++
			}
			if x != x0 {
				mmap[y][x0] = '.'
				mmap[y][x] = 'O'
			}
		}
	}
}

func plot(mmap [][]rune) {
	for y, row := range mmap {
		fmt.Printf("%2d: %s\n", len(mmap)-y, string(row))
	}
}

func rollCycle(mmap [][]rune) {
	rollUp(mmap)
	rollLeft(mmap)
	rollDown(mmap)
	rollRight(mmap)
}

func part2(lines []string) {
	timeStart := time.Now()
	mmap := parseMap(lines)

	visited := make(map[string]int)

	var i, loopStart int
	for i = 0; i < 1e9; i++ {
		key := fmt.Sprint(mmap)
		if step, ok := visited[key]; ok {
			fmt.Println("found loop at", i, "to step", step)
			loopStart = step
			break
		}
		visited[key] = i
		rollCycle(mmap)
	}
	left := (1e9 - loopStart) % (i - loopStart)
	fmt.Println("skipping to", 1e9-left)
	for ; left > 0; left-- {
		rollCycle(mmap)
	}

	sum := calcSum(mmap)
	fmt.Println("Part 2:", sum, "\tin", time.Since(timeStart))
}
