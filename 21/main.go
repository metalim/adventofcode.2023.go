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

type Pos struct {
	x, y int
}

func findStart(lines []string) Pos {
	for y, line := range lines {
		for x, char := range line {
			if char == 'S' {
				return Pos{x, y}
			}
		}
	}
	panic("no start found")
}

func part1(lines []string) {
	timeStart := time.Now()
	start := findStart(lines)
	next := []Pos{start}
	visited := map[Pos]bool{start: true}
	var cur []Pos
	for step := 0; step < 64; step++ {
		cur, next = next, cur[:0]
		for _, pos := range cur {
			for _, dir := range []Pos{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
				pos := Pos{pos.x + dir.x, pos.y + dir.y}
				if pos.x < 0 || pos.y < 0 || pos.x >= len(lines[0]) || pos.y >= len(lines) || visited[pos] {
					continue
				}
				visited[pos] = true
				switch lines[pos.y][pos.x] {
				case '.', 'S':
					next = append(next, pos)
				}
			}
		}
	}

	var count int
	for p := range visited {
		switch lines[p.y][p.x] {
		case '.', 'S':
			if (p.x+p.y)%2 == (start.x+start.y)%2 {
				count++
			}
		}
	}
	fmt.Println("Part 1:", count, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	start := findStart(lines)
	next := []Pos{start}
	visited := map[Pos]bool{start: true}
	var cur []Pos
	_, _, _ = cur, next, visited
	var count int
	fmt.Println("Part 2:", count, "\tin", time.Since(timeStart))
}
