package main

import (
	"fmt"
	"strings"
)

type Range struct {
	Min, Max int
}

const MaxInt = 1<<63 - 1

func part2ranges(lines []string) {
	seeds := toInts(strings.Fields(reSeeds.FindStringSubmatch(lines[0])[1]))
	maps := parseMaps(lines[1:])

	closest := MaxInt
	for i := 0; i < len(seeds); i += 2 {
		seed := seeds[i]
		length := seeds[i+1]
		pos := closestInRange(Range{seed, seed + length - 1}, maps)
		if closest > pos {
			closest = pos
		}
	}
	fmt.Println(closest)
}

func closestInRange(r Range, maps [][][]int) int {
	if len(maps) == 0 {
		return r.Min
	}

	if r.Min > r.Max {
		panic(fmt.Sprintf("invalid range: %v", r))
	}

	for _, row := range maps[0] {
		min := row[1]
		max := row[1] + row[2] - 1
		if r.Min > max || r.Max < min {
			continue
		}

		// range overlap
		delta := row[0] - row[1]
		if min < r.Min {
			min = r.Min
		}
		if max > r.Max {
			max = r.Max
		}
		closest := closestInRange(Range{min + delta, max + delta}, maps[1:])
		if r.Min < min {
			closestLeft := closestInRange(Range{r.Min, min - 1}, maps)
			if closest > closestLeft {
				closest = closestLeft
			}
		}
		if r.Max > max {
			closestRight := closestInRange(Range{max + 1, r.Max}, maps)
			if closest > closestRight {
				closest = closestRight
			}
		}
		return closest
	}

	// no matches, pass through to next layer
	return closestInRange(r, maps[1:])
}
