package main

import (
	"fmt"
	"os"
	"regexp"
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

type Step struct {
	name int
	cond string
	val  int
	out  string
}

type Workflow []Step
type Workflows map[string]Workflow
type Part [4]int
type Parts []Part

// px{a<2006:qkq,m>2090:A,rfg}
var reFlow = regexp.MustCompile(`^(\w+)\{(.*)\}$`)

// a<2006:qkq
var reRule = regexp.MustCompile(`^(\w+)([<>])(\d+):(\w+)$`)

// {x=787,m=2655,a=1222,s=2876}
var rePart = regexp.MustCompile(`^\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)\}$`)

func parse(lines []string) (Workflows, Parts) {
	workflows := Workflows{}
	parts := Parts{}

	var i int
	var line string
	for i, line = range lines {
		if line == "" {
			break
		}
		m := reFlow.FindStringSubmatch(line)
		if len(m) != 3 {
			panic("invalid line: " + line)
		}
		name := m[1]
		rules := strings.Split(m[2], ",")
		workflow := Workflow{}
		for _, rule := range rules {
			m2 := reRule.FindStringSubmatch(rule)
			if len(m2) != 5 {
				workflow = append(workflow, Step{name: -1, out: rule})
				continue
			}
			var valName int
			switch m2[1] {
			case "x":
				valName = 0
			case "m":
				valName = 1
			case "a":
				valName = 2
			case "s":
				valName = 3
			default:
				panic("invalid rule: " + rule)
			}
			val, err := strconv.Atoi(m2[3])
			catch(err)
			workflow = append(workflow, Step{
				name: valName,
				cond: m2[2],
				val:  val,
				out:  m2[4],
			})
		}
		workflows[name] = workflow
	}

	for _, line = range lines[i+1:] {
		m := rePart.FindStringSubmatch(line)
		if len(m) != 5 {
			panic("invalid line: " + line)
		}
		vals := toInts(m[1:])
		parts = append(parts, Part(vals))
	}

	return workflows, parts
}

func part1(lines []string) {
	timeStart := time.Now()
	workflows, parts := parse(lines)
	var sum int
	for _, part := range parts {
		state := "in"
		for state != "A" && state != "R" {
			flow, ok := workflows[state]
			if !ok {
				panic("invalid state: " + state)
			}
			for _, step := range flow {
				if step.name == -1 {
					state = step.out
					break
				}
				if step.cond == "<" && part[step.name] < step.val {
					state = step.out
					break
				}
				if step.cond == ">" && part[step.name] > step.val {
					state = step.out
					break
				}
			}
		}
		if state == "A" {
			sum += part[0] + part[1] + part[2] + part[3]
		}
	}
	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Println("Part 2:", "\tin", time.Since(timeStart))
}
