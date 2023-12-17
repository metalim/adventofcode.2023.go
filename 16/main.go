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

func traverse(lines []string, startBeam Beam) int {
	beams := []Beam{startBeam}
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
	return len(cells)
}

func part1(lines []string) {
	timeStart := time.Now()
	energized := traverse(lines, Beam{Pos{0, 0}, Right})
	fmt.Println("Part 1:", energized, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	var maxEnergized int
	for y := 0; y < len(lines); y++ {
		energized1 := traverse(lines, Beam{Pos{0, y}, Right})
		energized2 := traverse(lines, Beam{Pos{len(lines[0]) - 1, y}, Left})
		maxEnergized = max(maxEnergized, energized1, energized2)
	}
	for x := 0; x < len(lines[0]); x++ {
		energized1 := traverse(lines, Beam{Pos{x, 0}, Down})
		energized2 := traverse(lines, Beam{Pos{x, len(lines) - 1}, Up})
		maxEnergized = max(maxEnergized, energized1, energized2)
	}

	fmt.Println("Part 2:", maxEnergized, "\tin", time.Since(timeStart))
}
