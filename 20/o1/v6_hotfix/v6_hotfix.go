/*
➜ go run ./o1/ input1.txt
# github.com/metalim/adventofcode.2023.go/20/o1
o1/v5.go:304:13: declared and not used: stackIndex
o1/v5.go:409:22: cannot use out (variable of type []struct{name string; pulse PulseType}) as []pulseEvent value in argument to runPulses
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

/*
Описание изменений по сравнению с предыдущей версией:

- Удалён неиспользуемый stackIndex.
- Приведение out к []pulseEvent при вызове runPulses в findPeriodForSCC.
*/

type PulseType bool

const (
	Low  PulseType = false
	High PulseType = true
)

type ModuleType int

const (
	Normal ModuleType = iota
	FlipFlop
	Conjunction
	Broadcaster
)

type Module struct {
	mtype        ModuleType
	name         string
	dest         []string
	flipOn       bool
	conInputs    map[string]PulseType
	conInputList []string
}

func (m *Module) receive(from string, p PulseType) []struct {
	name  string
	pulse PulseType
} {
	switch m.mtype {
	case Broadcaster:
		var out []struct {
			name  string
			pulse PulseType
		}
		for _, d := range m.dest {
			out = append(out, struct {
				name  string
				pulse PulseType
			}{d, p})
		}
		return out
	case FlipFlop:
		if p == High {
			return nil
		}
		m.flipOn = !m.flipOn
		var outP PulseType
		if m.flipOn {
			outP = High
		} else {
			outP = Low
		}
		var out []struct {
			name  string
			pulse PulseType
		}
		for _, d := range m.dest {
			out = append(out, struct {
				name  string
				pulse PulseType
			}{d, outP})
		}
		return out
	case Conjunction:
		m.conInputs[from] = p
		allHigh := true
		for _, in := range m.conInputList {
			if m.conInputs[in] == Low {
				allHigh = false
				break
			}
		}
		var outP PulseType
		if allHigh {
			outP = Low
		} else {
			outP = High
		}
		var out []struct {
			name  string
			pulse PulseType
		}
		for _, d := range m.dest {
			out = append(out, struct {
				name  string
				pulse PulseType
			}{d, outP})
		}
		return out
	case Normal:
		var out []struct {
			name  string
			pulse PulseType
		}
		for _, d := range m.dest {
			out = append(out, struct {
				name  string
				pulse PulseType
			}{d, p})
		}
		return out
	}
	return nil
}

func reset(modules map[string]*Module) {
	for _, mod := range modules {
		mod.flipOn = false
		if mod.mtype == Conjunction {
			for _, in := range mod.conInputList {
				mod.conInputs[in] = Low
			}
		}
	}
}

type State struct {
	flipFlops map[string]bool
	conj      map[string][]PulseType
}

func getCurrentState(modules map[string]*Module) State {
	ff := make(map[string]bool)
	cj := make(map[string][]PulseType)
	for name, m := range modules {
		if m.mtype == FlipFlop {
			ff[name] = m.flipOn
		}
		if m.mtype == Conjunction {
			arr := make([]PulseType, len(m.conInputList))
			for i, inp := range m.conInputList {
				arr[i] = m.conInputs[inp]
			}
			cj[name] = arr
		}
	}
	return State{flipFlops: ff, conj: cj}
}

func applyState(modules map[string]*Module, s State) {
	for name, m := range modules {
		if m.mtype == FlipFlop {
			m.flipOn = s.flipFlops[name]
		}
		if m.mtype == Conjunction {
			for i, inp := range m.conInputList {
				m.conInputs[inp] = s.conj[m.name][i]
			}
		}
	}
}

func serializeState(s State) string {
	var sb strings.Builder
	sb.WriteString("FF:")
	ffnames := make([]string, 0, len(s.flipFlops))
	for n := range s.flipFlops {
		ffnames = append(ffnames, n)
	}
	sort.Strings(ffnames)
	for _, n := range ffnames {
		if s.flipFlops[n] {
			sb.WriteString(n + "=1,")
		} else {
			sb.WriteString(n + "=0,")
		}
	}
	sb.WriteString(";CONJ:")
	conjnames := make([]string, 0, len(s.conj))
	for n := range s.conj {
		conjnames = append(conjnames, n)
	}
	sort.Strings(conjnames)
	for _, n := range conjnames {
		sb.WriteString(n + "=")
		for _, v := range s.conj[n] {
			if v == High {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteString(",")
	}
	return sb.String()
}

type pulseEvent struct {
	from, to string
	ptype    PulseType
}

func runPulses(modules map[string]*Module, pulses []pulseEvent, countPulses *struct{ low, high int }, checkRx bool) bool {
	queue := pulses
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if p.ptype == Low {
			countPulses.low++
		} else {
			countPulses.high++
		}
		if checkRx && p.to == "rx" && p.ptype == Low {
			return true
		}
		mod, ok := modules[p.to]
		if !ok {
			continue
		}
		out := mod.receive(p.from, p.ptype)
		var newOut []pulseEvent
		for _, o := range out {
			newOut = append(newOut, pulseEvent{from: mod.name, to: o.name, ptype: o.pulse})
		}
		queue = append(queue, newOut...)
	}
	return false
}

func pressButton(modules map[string]*Module, countPulses *struct{ low, high int }, checkRx bool) bool {
	return runPulses(modules, []pulseEvent{{from: "button", to: "broadcaster", ptype: Low}}, countPulses, checkRx)
}

// Найдём все модули, влияющие на rx
func findModulesAffectingRx(modules map[string]*Module) map[string]bool {
	reverse := make(map[string][]string)
	for src, m := range modules {
		for _, d := range m.dest {
			reverse[d] = append(reverse[d], src)
		}
	}

	affected := make(map[string]bool)
	var stack []string
	stack = append(stack, "rx")
	affected["rx"] = true
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, p := range reverse[top] {
			if !affected[p] {
				affected[p] = true
				stack = append(stack, p)
			}
		}
	}
	return affected
}

// Алгоритм Тарьяна для SCC
func tarjanSCC(nodes []string, edges map[string][]string) [][]string {
	var index int
	stack := make([]string, 0)
	onStack := make(map[string]bool)
	indexes := make(map[string]int)
	lowlink := make(map[string]int)

	for _, n := range nodes {
		indexes[n] = -1
	}

	var sccs [][]string

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indexes[v] = index
		lowlink[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range edges[v] {
			if indexes[w] < 0 {
				strongConnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] {
				if indexes[w] < lowlink[v] {
					lowlink[v] = indexes[w]
				}
			}
		}

		if lowlink[v] == indexes[v] {
			var comp []string
			for {
				x := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[x] = false
				comp = append(comp, x)
				if x == v {
					break
				}
			}
			sccs = append(sccs, comp)
		}
	}

	for _, n := range nodes {
		if indexes[n] < 0 {
			strongConnect(n)
		}
	}

	return sccs
}

// Для упрощения: попытаемся найти период SCC
func findPeriodForSCC(modules map[string]*Module, scc []string) (int, int) {
	sccSet := make(map[string]bool)
	for _, n := range scc {
		sccSet[n] = true
	}

	initialState := getCurrentState(modules)
	seen := make(map[string]int)
	seen[serializeState(initialState)] = 0
	pressCount := 0

	testInputNode := scc[0]

	var dummyCount struct{ low, high int }
	for {
		pressCount++
		out := modules[testInputNode].receive("button", Low)
		var newOut []pulseEvent
		for _, o := range out {
			newOut = append(newOut, pulseEvent{from: modules[testInputNode].name, to: o.name, ptype: o.pulse})
		}
		runPulses(modules, newOut, &dummyCount, false)

		st := getCurrentState(modules)
		ser := serializeState(st)
		if oldI, ok := seen[ser]; ok {
			period := pressCount - oldI
			offset := oldI
			return period, offset
		}
		seen[ser] = pressCount
	}
}

// Решаем систему сравнений
func solveForRx(periods, offsets []int) int {
	gcd := func(x, y int) int {
		for y != 0 {
			x, y = y, x%y
		}
		return x
	}
	lcm := func(a, b int) int {
		return a / gcd(a, b) * b
	}

	allPeriod := 1
	for _, p := range periods {
		allPeriod = lcm(allPeriod, p)
		if allPeriod > 10000000 {
			return -1
		}
	}

	for n := 0; n <= allPeriod; n++ {
		match := true
		for i, p := range periods {
			if (n-offsets[i])%p != 0 {
				match = false
				break
			}
		}
		if match {
			return n
		}
	}
	return -1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <inputfile>")
		return
	}
	inputFile := os.Args[1]
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	modules := make(map[string]*Module)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, line := range lines {
		parts := strings.Split(line, "->")
		left := strings.TrimSpace(parts[0])
		right := ""
		if len(parts) > 1 {
			right = strings.TrimSpace(parts[1])
		}
		var dests []string
		if right != "" {
			for _, d := range strings.Split(right, ",") {
				dests = append(dests, strings.TrimSpace(d))
			}
		}
		var mtype ModuleType
		name := left
		if strings.HasPrefix(left, "%") {
			mtype = FlipFlop
			name = strings.TrimPrefix(left, "%")
		} else if strings.HasPrefix(left, "&") {
			mtype = Conjunction
			name = strings.TrimPrefix(left, "&")
		} else {
			if left == "broadcaster" {
				mtype = Broadcaster
			} else {
				mtype = Normal
			}
		}
		modules[name] = &Module{mtype: mtype, name: name, dest: dests, conInputs: make(map[string]PulseType)}
	}

	inputMap := make(map[string][]string)
	for src, mod := range modules {
		for _, d := range mod.dest {
			inputMap[d] = append(inputMap[d], src)
		}
	}
	for _, mod := range modules {
		if mod.mtype == Conjunction {
			mod.conInputList = inputMap[mod.name]
			for _, in := range mod.conInputList {
				mod.conInputs[in] = Low
			}
		}
	}

	// Part 1
	startPart1 := time.Now()
	reset(modules)
	var cp1 struct{ low, high int }
	for i := 0; i < 1000; i++ {
		pressButton(modules, &cp1, false)
	}
	part1Answer := int64(cp1.low) * int64(cp1.high)
	fmt.Println("Part 1 Answer:", part1Answer)
	fmt.Println("Time Part 1:", time.Since(startPart1))

	// Part 2
	startPart2 := time.Now()
	reset(modules)
	affected := findModulesAffectingRx(modules)

	filteredEdges := make(map[string][]string)
	var nodes []string
	for n := range modules {
		if affected[n] {
			nodes = append(nodes, n)
		}
	}
	for _, n := range nodes {
		var nd []string
		for _, d := range modules[n].dest {
			if affected[d] {
				nd = append(nd, d)
			}
		}
		filteredEdges[n] = nd
	}

	sccs := tarjanSCC(nodes, filteredEdges)

	var periods, offsets []int
	for _, c := range sccs {
		hasState := false
		for _, n := range c {
			m := modules[n]
			if m.mtype == FlipFlop || m.mtype == Conjunction {
				hasState = true
				break
			}
		}
		if hasState {
			save := getCurrentState(modules)
			p, off := findPeriodForSCC(modules, c)
			applyState(modules, save)
			periods = append(periods, p)
			offsets = append(offsets, off)
		}
	}

	var res int
	if len(periods) == 0 {
		reset(modules)
		var cp2 struct{ low, high int }
		foundRx := false
		for i := 1; i <= 100000; i++ {
			if pressButton(modules, &cp2, true) {
				foundRx = true
				res = i
				break
			}
		}
		if !foundRx {
			res = -1
		}
	} else {
		res = solveForRx(periods, offsets)
	}

	if res < 0 {
		fmt.Println("Part 2 Answer: no solution")
	} else {
		fmt.Println("Part 2 Answer:", res)
	}

	fmt.Println("Time Part 2:", time.Since(startPart2))
}
