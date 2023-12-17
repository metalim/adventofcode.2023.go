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

type Dir Pos

var (
	Up    = Dir{0, -1}
	Down  = Dir{0, 1}
	Left  = Dir{-1, 0}
	Right = Dir{1, 0}
)

type Beam struct {
	Pos
	dir Dir
}

func part1(lines []string) {
	timeStart := time.Now()

	grid := make([][]byte, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}

	beams := []Beam{{Pos{0, 0}, Right}}
	var next []Beam
	visited := map[Beam]bool{}
	cells := map[Pos]bool{}
	for len(beams) > 0 {
		for _, beam := range beams {
			if visited[beam] || beam.x < 0 || beam.x >= len(lines[0]) || beam.y < 0 || beam.y >= len(lines) {
				continue
			}
			visited[beam] = true
			cells[beam.Pos] = true
			c := lines[beam.y][beam.x]
			switch c {
			case '.':
				next = append(next, Beam{Pos{beam.x + beam.dir.x, beam.y + beam.dir.y}, beam.dir})
			case '\\':
				next = append(next, Beam{Pos{beam.x + beam.dir.y, beam.y + beam.dir.x}, Dir{beam.dir.y, beam.dir.x}})
			case '/':
				next = append(next, Beam{Pos{beam.x - beam.dir.y, beam.y - beam.dir.x}, Dir{-beam.dir.y, -beam.dir.x}})
			case '|':
				if beam.dir == Up || beam.dir == Down {
					// no effect
					next = append(next, Beam{Pos{beam.x, beam.y + beam.dir.y}, beam.dir})
				} else {
					// split to up and down
					next = append(next, Beam{Pos{beam.x, beam.y - 1}, Up})
					next = append(next, Beam{Pos{beam.x, beam.y + 1}, Down})
				}
			case '-':
				if beam.dir == Left || beam.dir == Right {
					// no effect
					next = append(next, Beam{Pos{beam.x + beam.dir.x, beam.y}, beam.dir})
				} else {
					// split to left and right
					next = append(next, Beam{Pos{beam.x - 1, beam.y}, Left})
					next = append(next, Beam{Pos{beam.x + 1, beam.y}, Right})
				}
			}
		}

		beams, next = next, beams[:0]
	}

	fmt.Println("Part 1:", len(cells), "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(timeStart))
}
