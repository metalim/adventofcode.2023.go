package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func part2math(lines []string) {
	timeStr := strings.Join(strings.Fields(lines[0])[1:], "")
	distanceStr := strings.Join(strings.Fields(lines[1])[1:], "")
	maxTime, err := strconv.Atoi(timeStr)
	catch(err)
	distance, err := strconv.Atoi(distanceStr)
	catch(err)

	// X*(maxTime-X) = distance
	// X*X - X*maxTime + distance = 0
	// X = (maxTime +- sqrt(maxTime*maxTime - 4*distance)) / 2
	sqrtD := math.Sqrt(float64(maxTime*maxTime - 4*distance))
	speedMin := int(math.Ceil((float64(maxTime) - sqrtD) / 2))
	speedMax := int(math.Floor((float64(maxTime) + sqrtD) / 2))
	fmt.Println("Part 2:", speedMax-speedMin+1)

}
