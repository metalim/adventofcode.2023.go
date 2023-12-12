package main

import (
	"fmt"
	"strings"
	"time"
)

func part2brute(lines []string) {
	timeStart := time.Now()
	seeds := toInts(strings.Fields(reSeeds.FindStringSubmatch(lines[0])[1]))
	maps := parseMaps(lines[1:])

	closest := -1
	ch := make(chan int, 100)
	for i := 0; i < len(seeds); i += 2 {
		go func(seed, length int) {
			closest := -1
			for j := 0; j < length; j++ {
				val := seed + j
				val = convert(val, maps)
				if closest > val || closest < 0 {
					closest = val
				}
			}
			ch <- closest
		}(seeds[i], seeds[i+1])
	}
	for i := 0; i < len(seeds); i += 2 {
		val := <-ch
		if closest > val || closest < 0 {
			closest = val
		}
	}
	fmt.Println("Part 2", closest, "\tin", time.Since(timeStart))
}
