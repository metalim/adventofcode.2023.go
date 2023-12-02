package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

var maxCubes = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
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

var (
	reGame = regexp.MustCompile(`(\d+): (.*)`)
	reTry  = regexp.MustCompile(`(\d+) ([^,;]+)`)
)

func part1(lines []string) {
	var idSum int

	for _, line := range lines {
		if line == "" {
			continue
		}

		game := reGame.FindStringSubmatch(line)
		if isValidGame(game[2]) {
			gameID, err := strconv.Atoi(game[1])
			catch(err)
			idSum += gameID
		}
	}
	fmt.Println(idSum)
}

func isValidGame(game string) bool {
	attempts := strings.Split(game, ";")
	for _, attempt := range attempts {
		cubes := reTry.FindAllStringSubmatch(attempt, -1)
		for _, cube := range cubes {
			n, err := strconv.Atoi(cube[1])
			catch(err)
			if maxCubes[cube[2]] < n {
				return false
			}
		}
	}
	return true
}

func part2(lines []string) {
	var sumPower int

	for _, line := range lines {
		if line == "" {
			continue
		}

		game := reGame.FindStringSubmatch(line)
		sumPower += getGamePower(game[2])
	}
	fmt.Println(sumPower)
}

func getGamePower(game string) int {
	minCubes := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	attempts := strings.Split(game, ";")
	for _, attempt := range attempts {
		cubes := reTry.FindAllStringSubmatch(attempt, -1)
		for _, cube := range cubes {
			n, err := strconv.Atoi(cube[1])
			catch(err)
			if minCubes[cube[2]] < n {
				minCubes[cube[2]] = n
			}
		}
	}

	pow := 1
	for _, v := range minCubes {
		pow *= v
	}
	return pow
}
