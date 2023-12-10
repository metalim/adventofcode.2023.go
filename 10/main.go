package main

import (
	"fmt"
	"os"
	"strings"
)

const PLOT = true

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
	// ┌─┐
	// │ │
	// └─┘
	switch char {
	case 'F', '┌':
		return []Pos{Down, Right}
	case '-', '─':
		return []Pos{Left, Right}
	case '7', '┐':
		return []Pos{Down, Left}
	case '|', '│':
		return []Pos{Up, Down}
	case 'L', '└':
		return []Pos{Up, Right}
	case 'J', '┘':
		return []Pos{Up, Left}
	case '.', ' ':
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
	if steps%2 != 0 {
		panic("Steps has to be even")
	}
	fmt.Println("Part 1:", steps/2)
}

func getLineChar(char rune) rune {
	// F-7
	// | |
	// L-J
	// ┌─┐
	// │ │
	// └─┘
	switch char {
	case 'F':
		return '┌'
	case '-':
		return '─'
	case '7':
		return '┐'
	case '|':
		return '│'
	case 'L':
		return '└'
	case 'J':
		return '┘'
	case '.':
		return ' '
	case 'S':
		return 'S'
	default:
		panic("Unknown char " + string(char))
	}
}

func createPath(lines []string, start Pos, startDirs []Pos) [][]rune {
	mmap := make([][]rune, len(lines))
	for y, line := range lines {
		mmap[y] = make([]rune, len(line))
		for x := range line {
			mmap[y][x] = ' '
		}
	}
	pos := start
	dir := startDirs[0]
	mmap[pos.y][pos.x] = 'S'
	for {
		pos = pos.Add(dir)
		if pos == start {
			break
		}
		mmap[pos.y][pos.x] = getLineChar(rune(lines[pos.y][pos.x]))
		dir = pos.GetNext(dir, lines)
	}
	return mmap
}

func part2(lines []string) {
	start := findStart(lines)

	// find connections
	startDirs := findStartDirections(lines, start)
	if len(startDirs) != 2 {
		panic("Start has to have 2 directions")
	}

	// create the path
	mmap := createPath(lines, start, startDirs)

	// mark left/right sides
	pos := start
	dir := startDirs[0]
	for {
		fill(mmap, pos.Add(dir.turnLeft()), ' ', '<')
		fill(mmap, pos.Add(dir.turnRight()), ' ', '>')
		pos = pos.Add(dir)
		fill(mmap, pos.Add(dir.turnLeft()), ' ', '<')
		fill(mmap, pos.Add(dir.turnRight()), ' ', '>')
		if pos == start {
			break
		}
		mmap[pos.y][pos.x] = getLineChar(rune(lines[pos.y][pos.x]))
		dir = pos.GetNext(dir, lines)
	}

	// count charInside
	charInside := '>'
	charOutside := '<'
	if mmap[0][0] == '>' {
		charInside = '<'
		charOutside = '>'
	}
	var inside, outside, nonMarked int
	for y, line := range mmap {
		for x, char := range line {
			switch char {
			case charInside:
				inside++
				mmap[y][x] = '#'
			case charOutside:
				outside++
				mmap[y][x] = ' '
			case ' ':
				nonMarked++
				mmap[y][x] = 'X'
			}
		}
	}
	plot(mmap)
	if nonMarked != 0 {
		panic(fmt.Sprint("Non marked:", nonMarked))
	}

	fmt.Println("Part 2:", inside)
}

func (dir Pos) turnLeft() Pos {
	switch dir {
	case Up:
		return Left
	case Left:
		return Down
	case Down:
		return Right
	case Right:
		return Up
	default:
		panic("Unknown dir")
	}
}

func (dir Pos) turnRight() Pos {
	switch dir {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	default:
		panic("Unknown dir")
	}
}

func fill(mmap [][]rune, pos Pos, from, to rune) {
	if pos.x < 0 || pos.x >= len(mmap[0]) || pos.y < 0 || pos.y >= len(mmap) {
		return
	}
	if mmap[pos.y][pos.x] != from {
		return
	}
	mmap[pos.y][pos.x] = to
	fill(mmap, pos.Add(Up), from, to)
	fill(mmap, pos.Add(Down), from, to)
	fill(mmap, pos.Add(Left), from, to)
	fill(mmap, pos.Add(Right), from, to)
}

func plot(mmap [][]rune) {
	if !PLOT {
		return
	}
	for _, line := range mmap {
		fmt.Println(string(line))
	}
	fmt.Println()
}
