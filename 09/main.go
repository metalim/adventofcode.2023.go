package main

import (
	"fmt"
	"os"
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

func next(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	subs, nonZero := getSubs(nums)
	if nonZero {
		return nums[len(nums)-1] + next(subs)
	}
	return nums[len(nums)-1]
}

func part1(lines []string) {
	var sum int
	for _, line := range lines {
		row := toInts(strings.Fields(line))
		sum += next(row)
	}

	fmt.Println("Part 1:", sum)
}

func getSubs(nums []int) (subs []int, nonZero bool) {
	subs = make([]int, len(nums)-1)
	for i, a := range nums[1:] {
		sub := a - nums[i]
		subs[i] = sub
		if sub != 0 {
			nonZero = true
		}
	}
	return subs, nonZero
}

func prev(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	subs, nonZero := getSubs(nums)
	if nonZero {
		return nums[0] - prev(subs)
	}
	return nums[0]
}

func part2(lines []string) {
	var sum int
	for _, line := range lines {
		row := toInts(strings.Fields(line))
		sum += prev(row)
	}

	fmt.Println("Part 2:", sum)
}
