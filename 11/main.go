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

type Pos struct {
	x, y int
}

func part1(lines []string) {
	expY := make([]int, len(lines))
	expX := make([]int, len(lines[0]))
	for x := range expX {
		expX[x] = 1
	}
	for y := range expY {
		expY[y] = 1
	}

	var galaxies []Pos
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				galaxies = append(galaxies, Pos{x, y})
				expX[x] = 0
				expY[y] = 0
			}
		}
	}
	for x := range expX[1:] {
		expX[x+1] += expX[x]
	}
	for y := range expY[1:] {
		expY[y+1] += expY[y]
	}
	for i, g := range galaxies {
		galaxies[i] = Pos{g.x + expX[g.x], g.y + expY[g.y]}
	}

	var sum int
	for i, g1 := range galaxies[1:] {
		for _, g2 := range galaxies[:i+1] {
			dist := abs(g1.x-g2.x) + abs(g1.y-g2.y)
			sum += dist
		}
	}
	fmt.Println("Part 1:", sum)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func part2(lines []string) {
	fmt.Println("Part 2:")
}