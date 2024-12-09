/*
первая часть выполнилась сразу, вторую часть я прервал после 32 минут ожидания

➜ go run ./o1/ input1.txt
Part 1 Answer: 861743850
Time Part 1: 8.622667ms

И кстати, serializeState не используется, используется только SerializeState — это другая функция.

Проанализируй наличие циклов в выводе независимых частей графа, и после их нахождения
ВЫЧИСЛИ минимальное количество требуемых нажатий. Повторю: ждать в брутфорсе получения
Low в rx бесполезно,т.к. число нажатий будет пятнадцатизначное
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

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
	for i := 0; i < len(ffnames); i++ {
		for j := i + 1; j < len(ffnames); j++ {
			if ffnames[j] < ffnames[i] {
				ffnames[i], ffnames[j] = ffnames[j], ffnames[i]
			}
		}
	}
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
	for i := 0; i < len(conjnames); i++ {
		for j := i + 1; j < len(conjnames); j++ {
			if conjnames[j] < conjnames[i] {
				conjnames[i], conjnames[j] = conjnames[j], conjnames[i]
			}
		}
	}
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

func runPulses(modules map[string]*Module, pulses []struct {
	from, to string
	ptype    PulseType
}, countPulses *struct{ low, high int }, checkRx bool) bool {
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
		for _, o := range out {
			queue = append(queue, struct {
				from, to string
				ptype    PulseType
			}{from: mod.name, to: o.name, ptype: o.pulse})
		}
	}
	return false
}

func pressButton(modules map[string]*Module, countPulses *struct{ low, high int }, checkRx bool) bool {
	return runPulses(modules, []struct {
		from, to string
		ptype    PulseType
	}{{from: "button", to: "broadcaster", ptype: Low}}, countPulses, checkRx)
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
	// Брютофорс долго, нужно искать цикл по состоянию.
	// Идея: симулируем нажатия пока:
	// 1. Не получим rx=low
	// 2. Или не обнаружим повтор состояния (цикл)
	// Если цикл без появления rx=low – значит rx никогда не будет low.

	startPart2 := time.Now()
	reset(modules)
	visited := make(map[string]int)
	steps := 0
	initialState := getCurrentState(modules)
	visited[serializeState(initialState)] = steps

	var cp2 struct{ low, high int }
	foundRx := false
	var answerSteps int

	for !foundRx {
		steps++
		if pressButton(modules, &cp2, true) {
			foundRx = true
			answerSteps = steps
			break
		}
		ns := getCurrentState(modules)
		key := serializeState(ns)
		if _, ok := visited[key]; ok {
			// Цикл обнаружен. Если rx не появился – никогда не появится.
			break
		} else {
			visited[key] = steps
		}
	}

	if foundRx {
		fmt.Println("Part 2 Answer:", answerSteps)
	} else {
		fmt.Println("Part 2 Answer: no solution (rx never gets low)")
	}
	fmt.Println("Time Part 2:", time.Since(startPart2))
}
