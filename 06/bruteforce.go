package main

import (
	"fmt"
	"strconv"
	"strings"
)

func part2brute(lines []string) {
	timeStr := strings.Join(strings.Fields(lines[0])[1:], "")
	distanceStr := strings.Join(strings.Fields(lines[1])[1:], "")
	maxTime, err := strconv.Atoi(timeStr)
	catch(err)
	distance, err := strconv.Atoi(distanceStr)
	catch(err)

	var won int
	for speed := 0; speed <= maxTime; speed++ {
		if speed*(maxTime-speed) > distance {
			won++
		}
	}

	fmt.Println("Part 2:", won)
}
