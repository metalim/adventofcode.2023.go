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

type pos struct {
	y, x int
}

func part1(lines []string) {
	used := map[pos]bool{}
	var sum int

	for y, line := range lines {
		for x, char := range line {
			if char == '.' || '0' <= char && char <= '9' {
				continue
			}

			for dy := -1; dy <= 1; dy++ {
				if y+dy < 0 || len(lines) <= y+dy {
					continue
				}
				for dx := -1; dx <= 1; dx++ {
					if dy == 0 && dx == 0 || x+dx < 0 || len(line) <= x+dx {
						continue
					}
					if used[pos{y + dy, x + dx}] {
						continue
					}

					c := lines[y+dy][x+dx]

					if c < '0' || '9' < c {
						continue
					}

					used[pos{y + dy, x + dx}] = true

					xl := x + dx
					for '0' <= c && c <= '9' {
						xl--
						if xl < 0 {
							break
						}
						c = lines[y+dy][xl]
						used[pos{y + dy, xl}] = true
					}
					xl++

					c = lines[y+dy][x+dx]
					xr := x + dx
					for '0' <= c && c <= '9' {
						xr++
						if len(line) <= xr {
							break
						}
						c = lines[y+dy][xr]
						used[pos{y + dy, xr}] = true
					}
					xr--

					n, err := strconv.Atoi(lines[y+dy][xl : xr+1])
					catch(err)

					sum += n
				}
			}
		}
	}
	fmt.Println(sum)
}

func part2(lines []string) {
	var sum int
	used := map[pos]bool{}

	for y, line := range lines {
		for x, char := range line {
			if char != '*' {
				continue
			}

			clear(used)
			mul := 1
			var numParts int

			for dy := -1; dy <= 1; dy++ {
				if y+dy < 0 || len(lines) <= y+dy {
					continue
				}
				for dx := -1; dx <= 1; dx++ {
					if dy == 0 && dx == 0 || x+dx < 0 || len(line) <= x+dx {
						continue
					}

					if used[pos{y + dy, x + dx}] {
						continue
					}

					c := lines[y+dy][x+dx]
					if c < '0' || '9' < c {
						continue
					}

					used[pos{y + dy, x + dx}] = true

					xl := x + dx
					for '0' <= c && c <= '9' {
						xl--
						if xl < 0 {
							break
						}
						c = lines[y+dy][xl]
						used[pos{y + dy, xl}] = true
					}
					xl++

					c = lines[y+dy][x+dx]
					xr := x + dx
					for '0' <= c && c <= '9' {
						xr++
						if len(line) <= xr {
							break
						}
						c = lines[y+dy][xr]
						used[pos{y + dy, xr}] = true
					}
					xr--

					n, err := strconv.Atoi(lines[y+dy][xl : xr+1])
					catch(err)

					mul *= n
					numParts++
				}
			}
			if numParts == 2 {
				sum += mul
			}
		}
	}
	fmt.Println(sum)
}
