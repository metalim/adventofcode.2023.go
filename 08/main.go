package main

import (
	"fmt"
	"os"
	"regexp"
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

var reRoute = regexp.MustCompile(`^(\w+) = \((\w+), (\w+)\)$`)

func parseRoutes(lines []string) map[string][2]string {
	routes := map[string][2]string{}
	for _, line := range lines {
		m := reRoute.FindStringSubmatch(line)
		if m == nil {
			panic("invalid route")
		}
		routes[m[1]] = [2]string{m[2], m[3]}
	}
	return routes
}

func part1(lines []string) {
	timeStart := time.Now()
	path := lines[0]
	routes := parseRoutes(lines[2:])

	var i, steps int
	next := "AAA"
	visited := map[string]bool{}
	for next != "ZZZ" {
		key := fmt.Sprintf("%s:%d", next, i)
		if visited[key] {
			fmt.Println("Part 1: loop", next, i, steps)
			return
		}
		visited[key] = true
		if _, ok := routes[next]; !ok {
			fmt.Println("Part 1: invalid", next, i, steps)
			return
		}
		switch path[i] {
		case 'L':
			next = routes[next][0]
		case 'R':
			next = routes[next][1]
		default:
			panic("unknown direction")
		}
		steps++
		i++
		if i >= len(path) {
			i = 0
		}
	}
	fmt.Println("Part 1:", steps, "\tin", time.Since(timeStart))
}

type Path struct {
	EndPos int
	EndVal string

	LoopStart    int
	LoopEnd      int
	LoopStartVal string

	// Cur
}

func part2(lines []string) {
	timeStart := time.Now()
	dirs := lines[0]
	routes := parseRoutes(lines[2:])

	var next []string
	for k := range routes {
		if k[2] == 'A' {
			next = append(next, k)
		}
	}

	paths := make([]Path, len(next))
	for i, n := range next {
		var j, steps int
		visited := map[string]int{}
		for {
			if n[2] == 'Z' && paths[i].EndPos == 0 {
				paths[i].EndPos = steps
				paths[i].EndVal = n
			}
			key := fmt.Sprintf("%s:%d", n, j)
			if _, ok := visited[key]; ok {
				paths[i].LoopStart = visited[key]
				paths[i].LoopStartVal = key
				paths[i].LoopEnd = steps
				break
			}
			visited[key] = steps
			if _, ok := routes[n]; !ok {
				fmt.Println("Part 2: invalid", n, i, j, steps)
				return
			}

			switch dirs[j] {
			case 'L':
				n = routes[n][0]
			case 'R':
				n = routes[n][1]
			default:
				panic("unknown direction")
			}

			steps++
			j++
			if j >= len(dirs) {
				j = 0
			}
		}
	}

	// AoC task is tuned to have simple solution, so we can just check two things:
	// 1. end position is in the loop
	for i, p := range paths {
		if p.EndPos < p.LoopStart || p.LoopEnd < p.LoopEnd {
			fmt.Println("Part 2: end position is not in loop", i, p.EndPos, p.LoopStart, p.LoopEnd)
			return
		}
	}
	// 2. end position is equal to loop length, so we can just use lcm
	for i, p := range paths[1:] {
		if p.EndPos != p.LoopEnd-p.LoopStart {
			fmt.Println("Part 2: end position is not equal to loop length", i, p.EndPos, p.LoopEnd-p.LoopStart)
			return
		}
	}

	nums := make([]int, len(paths))
	for i, p := range paths {
		nums[i] = p.EndPos
	}
	fmt.Println("Part 2:", lcm(nums...), "\tin", time.Since(timeStart))
}

func gcd(nums ...int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	a, b := nums[0], gcd(nums[1:]...)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(nums ...int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	a, b := nums[0], lcm(nums[1:]...)
	return a * b / gcd(a, b)
}
