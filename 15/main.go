package main

import (
	"fmt"
	"os"
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

	part1(lines)
	part2(lines)
}

func HASH(s string) int {
	var h int
	for _, c := range s {
		h = (h + int(c)) * 17 % 256
	}
	return h
}

func part1(lines []string) {
	timeStart := time.Now()
	var sum int
	for _, line := range lines {
		for _, s := range strings.Split(line, ",") {
			sum += HASH(s)
		}
	}

	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

type Lens struct {
	Label string
	Value int
}

type Hashmap [][]Lens

func (h Hashmap) Set(label string, value int) {
	i := HASH(label)
	for j, lens := range h[i] {
		if lens.Label == label {
			h[i][j].Value = value
			return
		}
	}
	h[i] = append(h[i], Lens{label, value})
}

func (h Hashmap) Delete(label string) {
	i := HASH(label)
	for j, lens := range h[i] {
		if lens.Label == label {
			h[i] = append(h[i][:j], h[i][j+1:]...)
			return
		}
	}
}

func part2(lines []string) {
	timeStart := time.Now()

	hashmap := make(Hashmap, 256)
	for _, line := range lines {
		for _, s := range strings.Split(line, ",") {
			if strings.HasSuffix(s, "-") {
				label := s[:len(s)-1]
				hashmap.Delete(label)
			} else {
				ss := strings.Split(s, "=")
				label := ss[0]
				value, err := strconv.Atoi(ss[1])
				catch(err)
				hashmap.Set(label, value)
			}
		}
	}

	var sum int
	for box, lens := range hashmap {
		for slot, l := range lens {
			sum += (box + 1) * (slot + 1) * l.Value
		}
	}
	fmt.Println("Part 2:", sum, "\tin", time.Since(timeStart))
}
