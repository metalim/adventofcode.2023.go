package main

import (
	"fmt"
	"os"
	"regexp"
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
	path := lines[0]
	routes := parseRoutes(lines[2:])

	var i, steps int
	next := "AAA"
	for next != "ZZZ" {
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
	fmt.Println("Part 1:", steps)
}

func part2(lines []string) {
}
