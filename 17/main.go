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

type (
	Dir  int
	Loss int
)

//go:generate stringer -type=Dir
const (
	Up Dir = iota
	Right
	Down
	Left
)

var Moves = []Pos{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

type Pos struct {
	x, y int
}
type PosDir struct {
	Pos
	dir   Dir
	steps int
}

type Crucible struct {
	PosDir
	loss Loss
}

func part1(lines []string) {
	timeStart := time.Now()
	minLoss := findMinLoss(lines, 0, 3)
	fmt.Println("Part 1:", minLoss, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	minLoss := findMinLoss(lines, 4, 10)
	fmt.Println("Part 2:", minLoss, "\tin", time.Since(timeStart))
}

func findMinLoss(lines []string, minSteps, maxSteps int) Loss {
	visited := make(map[PosDir]Loss)

	next := []Crucible{{PosDir{Pos{0, 0}, Right, 0}, 0}, {PosDir{Pos{0, 0}, Down, 0}, 0}}
	var current []Crucible
	var exit bool
	for !exit {
		exit = true
		current, next = next, current[:0]
		for _, cruc := range current {
			if v, ok := visited[cruc.PosDir]; ok {
				if v <= cruc.loss {
					continue
				}
			}
			exit = false
			visited[cruc.PosDir] = cruc.loss

			for turn := -1; turn <= 1; turn += 1 {
				newDir := (cruc.dir + Dir(turn) + 4) % 4
				steps := cruc.steps
				if turn == 0 {
					if cruc.steps >= maxSteps {
						continue
					}
				} else {
					if cruc.steps < minSteps {
						continue
					}
					steps = 0
				}
				newPosDir := PosDir{Pos{cruc.x + Moves[newDir].x, cruc.y + Moves[newDir].y}, newDir, steps + 1}
				if newPosDir.x < 0 || newPosDir.x >= len(lines[0]) || newPosDir.y < 0 || newPosDir.y >= len(lines) {
					continue
				}
				cruc := Crucible{newPosDir, cruc.loss + Loss(lines[newPosDir.y][newPosDir.x]-'0')}
				next = append(next, cruc)
			}
		}
	}

	var minLoss Loss = 1<<63 - 1
	for k, v := range visited {
		if k.x == len(lines[0])-1 && k.y == len(lines)-1 {
			// fmt.Printf("pos: %v, dir: %s, step: %d, loss: %v\n", k.Pos, k.dir, k.steps, v)
			if minLoss > v {
				minLoss = v
			}
		}
	}
	return minLoss
}
