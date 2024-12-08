package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const Part1Presses = 1000
const Part2Exit = "rx"

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	exit := flag.String("exit", Part2Exit, "exit node")
	export := flag.Bool("export", false, "export dot file")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(0)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	input := parseInput(string(bs))
	if *export {
		exportDot(input, flag.Arg(0)+".dot")
	}

	part1(input)
	part2(input, *exit)
}

const (
	Broadcaster = ""
	FlipFlop    = "%"
	Conjunction = "&"
	Untyped     = "?"

	Low  Level = false
	High Level = true
)

type Level bool

func (l Level) String() string {
	if l {
		return "High"
	}
	return "Low"
}

type Node struct {
	Type string
	Name string
	Dest []string

	state  Level            // for FlipFlop nodes only
	inputs map[string]Level // for Conjunction nodes only
}

func NewNode(typ, name string, dest []string) *Node {
	return &Node{Type: typ, Name: name, Dest: dest, inputs: make(map[string]Level)}
}

var reNode = regexp.MustCompile(`^([%&]?)(\w+) -> (.*)$`)

func parseInput(input string) map[string]*Node {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	nodes := make(map[string]*Node)
	for _, line := range lines {
		m := reNode.FindStringSubmatch(line)
		if len(m) == 0 {
			panic(fmt.Sprintf("Invalid line: %s", line))
		}
		nodes[m[2]] = NewNode(m[1], m[2], strings.Split(m[3], ", "))
	}

	return nodes
}

type Pulse struct {
	source, dest string
	level        Level
}

type Processor struct {
	nodes         map[string]*Node
	pulses        []Pulse
	lowCount      int
	highCount     int
	buttonPresses int

	monitors map[string]Monitor
}

func NewProcessor(nodes map[string]*Node) *Processor {
	for _, node := range nodes {
		node.state = Low
		for _, dest := range node.Dest {
			if _, ok := nodes[dest]; !ok {
				nodes[dest] = NewNode(Untyped, dest, nil)
			}
			nodes[dest].inputs[node.Name] = Low
		}
	}

	return &Processor{nodes: nodes, monitors: make(map[string]Monitor)}
}

func (s *Processor) Send(source, dest string, level Level) {
	s.pulses = append(s.pulses, Pulse{source: source, dest: dest, level: level})
	switch level {
	case High:
		s.highCount++
	case Low:
		s.lowCount++
	default:
		panic(fmt.Sprintf("Unknown level: %v", level))
	}
}

func (s *Processor) PressButton() {
	s.buttonPresses++
	s.Send("button", "broadcaster", Low)

	for i := 0; i < len(s.pulses); i++ {
		pulse := s.pulses[i]
		if monitor, ok := s.monitors[pulse.dest]; ok {
			if pulse.level == monitor.level {
				monitor.fn(pulse)
			}
		}
		node := s.nodes[pulse.dest]
		switch node.Type {
		case Broadcaster:
			for _, dest := range node.Dest {
				s.Send(node.Name, dest, pulse.level)
			}
		case FlipFlop:
			if pulse.level == Low {
				switch node.state {
				case Low:
					node.state = High
				case High:
					node.state = Low
				default:
					panic(fmt.Sprintf("Unknown state: %v", node.state))
				}
				for _, dest := range node.Dest {
					s.Send(node.Name, dest, node.state)
				}
			}
		case Conjunction:
			node.inputs[pulse.source] = pulse.level
			output := Low // output is Low, if all inputs are High
			for _, input := range node.inputs {
				if input == Low {
					output = High
					break
				}
			}
			for _, dest := range node.Dest {
				s.Send(node.Name, dest, output)
			}
		case Untyped:
			// ignore
		default:
			panic(fmt.Sprintf("Unknown node type: %s", node.Type))
		}
	}
	s.pulses = s.pulses[:0]
}

type MonitorFunc func(pulse Pulse)

type Monitor struct {
	name  string
	level Level
	fn    MonitorFunc
}

func (s *Processor) AddMonitor(name string, level Level, fn MonitorFunc) {
	s.monitors[name] = Monitor{name: name, level: level, fn: fn}
}

func part1(nodes map[string]*Node) {
	timeStart := time.Now()
	s := NewProcessor(nodes)
	for i := 0; i < Part1Presses; i++ {
		s.PressButton()
	}
	fmt.Printf("Low: %d, High: %d\n", s.lowCount, s.highCount)
	fmt.Printf("Part 1: %d\t\tin %v\n", s.lowCount*s.highCount, time.Since(timeStart))
}

func part2(nodes map[string]*Node, exit string) {
	timeStart := time.Now()

	s := NewProcessor(nodes)
	if _, ok := nodes[exit]; !ok {
		Errorf("Exit node %q not found\n", exit)
	}
	if len(nodes[exit].inputs) != 1 {
		Errorf("Exit node %q has %d inputs, expected 1\n", exit, len(nodes[exit].inputs))
	}
	var periods []int
	var expected int
	for preExit := range nodes[exit].inputs {
		if nodes[preExit].Type != Conjunction {
			Errorf("Node %q is not a Conjunction\n", preExit)
		}
		for period := range nodes[preExit].inputs {
			expected++
			s.AddMonitor(period, Low, func(pulse Pulse) {
				fmt.Printf("Press %d: %s %v -> %s\n", s.buttonPresses, pulse.source, pulse.level, pulse.dest)
				periods = append(periods, s.buttonPresses)
				expected--
			})
		}
	}

	for expected > 0 {
		s.PressButton()
	}

	mul := 1
	for _, p := range periods {
		mul *= p
	}
	fmt.Printf("Part 2: %d\t\tin %v\n", mul, time.Since(timeStart))
}

func Errorf(format string, args ...any) {
	fmt.Printf("Error: "+format, args...)
	os.Exit(1)
}
