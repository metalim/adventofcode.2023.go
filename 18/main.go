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

	part1Brute(lines)
	part1Smart(lines)
	part2(lines)
}

type Pos struct {
	x, y int
}

// use Shoelace formula
func part1Smart(lines []string) {
	timeStart := time.Now()
	var pos Pos
	var sum int
	var pathLen int
	for _, line := range lines {
		ss := strings.Fields(line)
		dir := ss[0]
		l, err := strconv.Atoi(ss[1])
		catch(err)

		pos0 := pos
		switch dir {
		case "R":
			pos.x += l
		case "D":
			pos.y += l
		case "L":
			pos.x -= l
		case "U":
			pos.y -= l
		}
		pathLen += l
		sum += pos.x*pos0.y - pos0.x*pos.y
	}

	fmt.Println("Part 1:", (abs(sum)+pathLen)/2+1, "\tin", time.Since(timeStart))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var reColor = regexp.MustCompile(`\(#(.*)(.)\)`)

func part2(lines []string) {
	timeStart := time.Now()
	var pos Pos
	var sum, pathLen int
	for _, line := range lines {
		m := reColor.FindStringSubmatch(line)
		dist64, err := strconv.ParseInt(m[1], 16, 64)
		catch(err)
		dir := m[2]
		dist := int(dist64)

		pos0 := pos
		switch dir {
		case "0":
			pos.x += dist
		case "1":
			pos.y += dist
		case "2":
			pos.x -= dist
		case "3":
			pos.y -= dist
		}
		pathLen += dist
		sum += pos.x*pos0.y - pos0.x*pos.y
	}

	fmt.Println("Part 2:", (abs(sum)+pathLen)/2+1, "\tin", time.Since(timeStart))
}
