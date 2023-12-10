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

func (p Pos) Add(other Pos) Pos {
	return Pos{x: p.x + other.x, y: p.y + other.y}
}

func findStart(lines []string) Pos {
	for y, line := range lines {
		for x, char := range line {
			if char == 'S' {
				return Pos{x: x, y: y}
			}
		}
	}
	panic("No start found")
}

var (
	Up    = Pos{x: 0, y: -1}
	Down  = Pos{x: 0, y: 1}
	Left  = Pos{x: -1, y: 0}
	Right = Pos{x: 1, y: 0}
)

func getDirs(char rune) []Pos {
	switch char {
	case 'J':
		return []Pos{Up, Left}
	case 'L':
		return []Pos{Up, Right}
	case '7':
		return []Pos{Down, Left}
	case 'F':
		return []Pos{Down, Right}
	case '|':
		return []Pos{Up, Down}
	case '-':
		return []Pos{Left, Right}
	case '.':
		return nil
	default:
		panic("Unknown char " + string(char))
	}
}

func findStartDirections(lines []string, start Pos) []Pos {
	startDirs := make([]Pos, 0, 2)

	for _, dir := range []Pos{Up, Down, Left, Right} {
		pos := start.Add(dir)
		if pos.x < 0 || pos.x >= len(lines[0]) || pos.y < 0 || pos.y >= len(lines) {
			continue
		}
		dirs := getDirs(rune(lines[pos.y][pos.x]))
		if dirs == nil {
			continue
		}
		for _, dir2 := range dirs {
			if pos.Add(dir2) == start {
				startDirs = append(startDirs, dir)
				break
			}
		}
	}
	return startDirs
}

func (pos Pos) GetNext(moveDir Pos, lines []string) Pos {
	dirs := getDirs(rune(lines[pos.y][pos.x]))
	if len(dirs) != 2 {
		panic("Expected 2 directions")
	}
	for _, dir := range dirs {
		if pos.Add(dir).Add(moveDir) != pos {
			return dir
		}
	}
	panic("No next direction found")
}

func part1(lines []string) {
	start := findStart(lines)

	// find connections
	startDirs := findStartDirections(lines, start)
	if len(startDirs) != 2 {
		panic("Start has to have 2 directions")
	}

	var steps int
	pos := start
	dir := startDirs[0]
	for {
		pos = pos.Add(dir)
		steps++
		if pos == start {
			break
		}
		dir = pos.GetNext(dir, lines)
	}
	fmt.Println("Part 1:", (steps)/2)
}

func part2(lines []string) {
	fmt.Println("Part 2:")
}
