package main

import (
	"fmt"
	"os"
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

type Pos struct {
	x, y int
}

func part1(lines []string) {
	timeStart := time.Now()
	grid := map[Pos]byte{}
	var pos Pos
	for _, line := range lines {
		ss := strings.Fields(line)
		dir := ss[0]
		l, err := strconv.Atoi(ss[1])
		catch(err)

		for ; l > 0; l-- {
			switch dir {
			case "U":
				pos.y--
			case "D":
				pos.y++
			case "L":
				pos.x--
			case "R":
				pos.x++
			}
			grid[pos] = '#'
		}
	}
	rect := findRect(grid)
	fillOuter(grid, rect)
	inner := countInner(grid, rect)
	// plot(grid, rect)
	fmt.Println("Part 1:", inner, "\tin", time.Since(timeStart))
}

type Rect struct {
	minX, minY, maxX, maxY int
}

func findRect(grid map[Pos]byte) Rect {
	var out Rect
	for pos := range grid {
		out.minX = min(out.minX, pos.x)
		out.minY = min(out.minY, pos.y)
		out.maxX = max(out.maxX, pos.x)
		out.maxY = max(out.maxY, pos.y)
	}
	return out
}

func fillOuter(grid map[Pos]byte, rect Rect) {
	var next []Pos
	var cur []Pos
	for y := rect.minY; y <= rect.maxY; y++ {
		cur = append(cur, Pos{rect.minX, y})
		cur = append(cur, Pos{rect.maxX, y})
	}
	for x := rect.minX; x <= rect.maxX; x++ {
		cur = append(cur, Pos{x, rect.minY})
		cur = append(cur, Pos{x, rect.maxY})
	}
	for len(cur) > 0 {
		for _, pos := range cur {
			if pos.x < rect.minX || pos.x > rect.maxX || pos.y < rect.minY || pos.y > rect.maxY {
				continue
			}
			if _, ok := grid[pos]; ok {
				continue
			}
			grid[pos] = '.'
			next = append(next, Pos{pos.x - 1, pos.y})
			next = append(next, Pos{pos.x + 1, pos.y})
			next = append(next, Pos{pos.x, pos.y - 1})
			next = append(next, Pos{pos.x, pos.y + 1})
		}
		cur, next = next, cur[:0]
	}
}

func countInner(grid map[Pos]byte, rect Rect) int {
	var count int
	for y := rect.minY; y <= rect.maxY; y++ {
		for x := rect.minX; x <= rect.maxX; x++ {
			if grid[Pos{x, y}] != '.' {
				count++
			}
		}
	}
	return count
}
func plot(grid map[Pos]byte, rect Rect) {
	for y := rect.minY; y <= rect.maxY; y++ {
		for x := rect.minX; x <= rect.maxX; x++ {
			if v, ok := grid[Pos{x, y}]; ok {
				fmt.Print(string(v))
			} else {
				fmt.Print("x")
			}
		}
		fmt.Println()
	}
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(timeStart))
}
