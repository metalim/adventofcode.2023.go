package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func toInts(ss []string) []int {
	is := make([]int, len(ss))
	for i, s := range ss {
		n, err := strconv.Atoi(s)
		catch(err)
		is[i] = n
	}
	return is
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
	// part2brute(lines)
	part2math(lines)
}

func part1(lines []string) {
	times := toInts(strings.Fields(lines[0])[1:])
	distances := toInts(strings.Fields(lines[1])[1:])
	mul := 1
	for i, maxTime := range times {
		var won int
		for speed := 0; speed <= maxTime; speed++ {
			if speed*(maxTime-speed) > distances[i] {
				won++
			}
		}
		mul *= won
	}
	fmt.Println("Part 1:", mul)
}
