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
	part2(lines)
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

func part2(lines []string) {
	timeStr := strings.Join(strings.Fields(lines[0])[1:], "")
	distanceStr := strings.Join(strings.Fields(lines[1])[1:], "")
	maxTime, err := strconv.Atoi(timeStr)
	catch(err)
	distance, err := strconv.Atoi(distanceStr)
	catch(err)

	// we can do binary search here, but since T is just 10s of millions, we do brute force
	var won int
	for speed := 0; speed <= maxTime; speed++ {
		if speed*(maxTime-speed) > distance {
			won++
		}
	}

	fmt.Println("Part 2:", won)
}
