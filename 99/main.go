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

func part1(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func part2(lines []string) {
}
