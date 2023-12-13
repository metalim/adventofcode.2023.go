package main

import (
	"fmt"
	"os"
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

	part0(lines)
	part1(lines)
	part2(lines)
}

func getCounts(field string) (damaged, unknown int) {
	for _, c := range field {
		switch c {
		case '#':
			damaged++
		case '?':
			unknown++
		case '.':
			// Do nothing
		default:
			panic("Unknown field: " + string(c))
		}
	}
	return damaged, unknown
}

func getSum(seq []int) (sum int) {
	for _, n := range seq {
		sum += n
	}
	return sum
}

func getVariants(field string, seq []int) (out int) {
	fmt.Printf("%q %v\n", field, seq)
	defer func() {
		fmt.Printf("%q %v -> %d\n", field, seq, out)
	}()
	if len(field) == 0 && len(seq) == 0 {
		return 1
	}
	if len(field) == 0 {
		return 0
	}

	seqSum := getSum(seq)
	if len(field) < seqSum+len(seq)-1 {
		return 0
	}
	damaged, unknown := getCounts(field)
	if seqSum < damaged || damaged+unknown < seqSum {
		return 0
	}

	var sum int
	switch field[0] {
	case '.':
		sum += getVariantsMemoized(field[1:], seq)

	case '?':
		fmt.Println("1. treat ? as .")
		sum += getVariantsMemoized(field[1:], seq)
		fmt.Println("2. treat ? as #")
		fallthrough

	case '#':
		if seq[0] == 1 {
			// next should be . or ?
			if len(seq) == 1 {
				sum += getVariantsMemoized(field[1:], seq[1:])
			} else {
				if field[1] == '.' || field[1] == '?' {
					sum += getVariantsMemoized(field[2:], seq[1:])
				}
			}
		} else {
			seq[0]--
			sum += getVariantsMemoized(field[1:], seq)
			seq[0]++
		}

	default:
		panic("Unknown field: " + string(field[0]))
	}
	return sum
}

type Args struct {
	field string
	seq   string
}

var memo = map[Args]int{}

func getVariantsMemoized(field string, seq []int) int {
	seqStr := fmt.Sprint(seq)

	if memoized, ok := memo[Args{field, seqStr}]; ok {
		fmt.Printf("%q %v -> %d [memo]\n", field, seq, memoized)
		return memoized
	}
	result := getVariants(field, seq)
	memo[Args{field, seqStr}] = result
	return result
}

func part0(lines []string) {
	// get
}

func part1(lines []string) {
	start := time.Now()
	var sum int
	for _, line := range lines {
		fields := strings.Fields(line)
		nums := toInts(strings.Split(fields[1], ","))
		vars := getVariantsMemoized(fields[0], nums)
		fmt.Println(line, vars)
		sum += vars
	}

	fmt.Println("Part 1:", sum, "\tin", time.Since(start))
}

func part2(lines []string) {
	start := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(start))
}
