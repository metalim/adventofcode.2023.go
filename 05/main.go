package main

import (
	"fmt"
	"os"
	"regexp"
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
	// part2brute(lines) // takes 1m30s
	part2ranges(lines)
}

var reSeeds = regexp.MustCompile(`^seeds: (.*)$`)

func toInts(ss []string) []int {
	is := make([]int, len(ss))
	for i, s := range ss {
		n, err := strconv.Atoi(s)
		catch(err)
		is[i] = n
	}
	return is
}

func parseMaps(lines []string) [][][]int {
	maps := make([][][]int, 0, 7) // map, row, numbers

	var skipNext bool
	var mmap [][]int
	for _, line := range lines {
		if skipNext {
			skipNext = false
			continue
		}
		if line == "" {
			skipNext = true
			if mmap != nil {
				maps = append(maps, mmap)
			}
			mmap = nil
			continue
		}

		mmap = append(mmap, toInts(strings.Fields(line)))
	}
	if mmap != nil {
		maps = append(maps, mmap)
	}

	return maps
}

func convert(val int, maps [][][]int) int {
	for _, mmap := range maps {
		for _, row := range mmap {
			if val >= row[1] && val < row[1]+row[2] {
				val = val - row[1] + row[0]
				break
			}
		}
	}
	return val
}

func part1(lines []string) {
	timeStart := time.Now()
	strSeeds := reSeeds.FindStringSubmatch(lines[0])[1]
	seeds := toInts(strings.Fields(strSeeds))

	maps := parseMaps(lines[1:])

	closest := -1
	for _, val := range seeds {
		val = convert(val, maps)
		if closest > val || closest < 0 {
			closest = val
		}
	}
	fmt.Println("Part 1:", closest, "\tin", time.Since(timeStart))
}
