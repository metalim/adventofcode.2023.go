package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
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

var reBrick = regexp.MustCompile(`^(\d+),(\d+),(\d+)~(\d+),(\d+),(\d+)$`)

type Brick struct {
	pos         [6]int
	supportedBy []*Brick
	supports    []*Brick
}

func parse(lines []string) []*Brick {
	bricks := make([]*Brick, 0, len(lines))
	for _, line := range lines {
		matches := reBrick.FindStringSubmatch(line)
		if matches == nil {
			panic("invalid brick")
		}
		brick := &Brick{
			pos: [6]int(toInts(matches[1:])),
		}
		bricks = append(bricks, brick)
	}
	return bricks
}

func (b Brick) Volume() int {
	return (b.pos[3] - b.pos[0] + 1) *
		(b.pos[4] - b.pos[1] + 1) *
		(b.pos[5] - b.pos[2] + 1)
}

func (b Brick) Dimensions() [3]int {
	return [3]int{
		b.pos[3] - b.pos[0] + 1,
		b.pos[4] - b.pos[1] + 1,
		b.pos[5] - b.pos[2] + 1,
	}
}

func (b Brick) Intersect(bricks []*Brick) (collisions []*Brick) {
	for _, brick := range bricks {
		if b.pos[0] <= brick.pos[3] && b.pos[3] >= brick.pos[0] &&
			b.pos[1] <= brick.pos[4] && b.pos[4] >= brick.pos[1] &&
			b.pos[2] <= brick.pos[5] && b.pos[5] >= brick.pos[2] {
			collisions = append(collisions, brick)
		}
	}
	return collisions
}

func (b *Brick) Move(dx, dy, dz int) {
	b.pos[0] += dx
	b.pos[1] += dy
	b.pos[2] += dz
	b.pos[3] += dx
	b.pos[4] += dy
	b.pos[5] += dz
}

func printBricks(bricks []*Brick) {
	for _, brick := range bricks {
		fmt.Printf("%v: %v, %v\n", brick.pos, brick.supportedBy, brick.supports)
	}
	fmt.Println()
}

func part1(lines []string) {
	timeStart := time.Now()
	bricks := parse(lines)
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].pos[5] < bricks[j].pos[5]
	})

	// printBricks(bricks)

	for i, brick := range bricks {
		for brick.pos[5] > 0 {
			collisions := brick.Intersect(bricks[:i])
			if len(collisions) == 0 {
				brick.Move(0, 0, -1)
				continue
			}
			brick.supportedBy = collisions
			for _, collision := range collisions {
				collision.supports = append(collision.supports, brick)
			}
			break
		}
		brick.Move(0, 0, 1)
	}

	printBricks(bricks)

	var canBeRemoved int
	for _, brick := range bricks {
		redundant := true
		for _, above := range brick.supports {
			if len(above.supportedBy) < 2 {
				redundant = false
				break
			}
		}
		if redundant {
			canBeRemoved++
		}
	}
	fmt.Printf("Part 1: %d of %d \tin %v\n", canBeRemoved, len(bricks), time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Printf("Part 2: \tin %v\n", time.Since(timeStart))
}
