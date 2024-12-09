/*
```
➜ go run ./o1/v11 input1.txt
No subsystems found in input
```
Эй! не ломай достигнутое!
ВОЗЬМИ КОД РАБОЧЕЙ ПРОГРАММЫ и добавь туда вычисление результата. НЕ ЛОМАЙ!!!!!!
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*
Берём код из предыдущей рабочей версии (v10), где программа успешно находила периоды и смещения.

Изменения:
- В момент, когда мы вычисляем период и смещение для каждой подсистемы, будем сохранять их в массив.
- После того как все подсистемы обработаны, посчитаем итоговый результат:
  Предполагается, что все offset равны (в примере offset=1).
  Тогда результат: n = offset + LCM(all periods).

Добавляем в конце вычисление этого n и выводим.

Ничего другого не ломаем.
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

type pulseEvent struct {
	from, to string
	ptype    PulseType
}

func runAllPulses(modules map[string]*Module, pulses []pulseEvent) {
	queue := pulses
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		mod, ok := modules[p.to]
		if !ok {
			continue
		}
		out := mod.receive(p.from, p.ptype)
		for _, o := range out {
			queue = append(queue, pulseEvent{from: mod.name, to: o.name, ptype: o.pulse})
		}
	}
}

func pressButton(modules map[string]*Module) {
	runAllPulses(modules, []pulseEvent{{from: "button", to: "broadcaster", ptype: Low}})
}

func findBroadcaster(modules map[string]*Module) string {
	for n, m := range modules {
		if m.mtype == Broadcaster {
			return n
		}
	}
	return ""
}

func findAggregator(modules map[string]*Module) string {
	for n, m := range modules {
		for _, d := range m.dest {
			if d == "rx" && m.mtype == Conjunction {
				return n
			}
		}
	}
	return ""
}

func buildReverseGraph(modules map[string]*Module) map[string][]string {
	rg := make(map[string][]string)
	for src, m := range modules {
		for _, d := range m.dest {
			rg[d] = append(rg[d], src)
		}
	}
	return rg
}

func subsystemState(modules map[string]*Module, subNodes map[string]bool) string {
	var ffNames []string
	var conjNames []string
	ffMap := make(map[string]bool)
	conjMap := make(map[string][]PulseType)

	for n := range subNodes {
		m := modules[n]
		if m.mtype == FlipFlop {
			ffMap[n] = m.flipOn
			ffNames = append(ffNames, n)
		}
		if m.mtype == Conjunction {
			conjMap[n] = make([]PulseType, len(m.conInputList))
			for i, inp := range m.conInputList {
				conjMap[n][i] = m.conInputs[inp]
			}
			conjNames = append(conjNames, n)
		}
	}
	sort.Strings(ffNames)
	sort.Strings(conjNames)

	sb := strings.Builder{}
	sb.WriteString("FF:")
	for _, n := range ffNames {
		if ffMap[n] {
			sb.WriteString(n + "=1,")
		} else {
			sb.WriteString(n + "=0,")
		}
	}
	sb.WriteString(";CONJ:")
	for _, n := range conjNames {
		sb.WriteString(n + "=")
		for _, v := range conjMap[n] {
			if v == High {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteByte(',')
	}
	return sb.String()
}

func findSubsystemNodes(modules map[string]*Module, reverseGraph map[string][]string, broadcaster string, startNode string) map[string]bool {
	subNodes := make(map[string]bool)
	broadDest := make(map[string]bool)
	for _, d := range modules[broadcaster].dest {
		broadDest[d] = true
	}

	visited := make(map[string]bool)
	var stack []string
	stack = append(stack, startNode)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[cur] {
			continue
		}
		visited[cur] = true
		subNodes[cur] = true

		if cur == broadcaster || broadDest[cur] {
			continue
		}
		for _, p := range reverseGraph[cur] {
			stack = append(stack, p)
		}
	}
	return subNodes
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
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
			sort.Strings(mod.conInputList)
			for _, in := range mod.conInputList {
				mod.conInputs[in] = Low
			}
		}
	}

	broad := findBroadcaster(modules)
	if broad == "" {
		fmt.Println("No broadcaster found")
		return
	}

	aggregator := findAggregator(modules)
	if aggregator == "" {
		fmt.Println("No aggregator found")
		return
	}
	aggMod := modules[aggregator]
	if aggMod.mtype != Conjunction {
		fmt.Println("Aggregator is not a conjunction")
		return
	}
	rxFound := false
	for _, d := range aggMod.dest {
		if d == "rx" {
			rxFound = true
			break
		}
	}
	if !rxFound {
		fmt.Println("Aggregator does not lead to rx")
		return
	}
	if len(aggMod.conInputList) == 0 {
		fmt.Println("Aggregator has no inputs")
		return
	}

	reverseGraph := buildReverseGraph(modules)

	var results []struct {
		inm    string
		period int
		offset int
	}

	// Периоды и смещения
	for _, inm := range aggMod.conInputList {
		subNodes := findSubsystemNodes(modules, reverseGraph, broad, inm)

		reset(modules)
		stateMap := make(map[string]int)
		step := 0
		st := subsystemState(modules, subNodes)
		stateMap[st] = step

		for {
			step++
			pressButton(modules)
			st2 := subsystemState(modules, subNodes)
			if oldStep, ok := stateMap[st2]; ok {
				offset := oldStep
				period := step - oldStep
				fmt.Printf("Subsystem input %s: period=%d, offset=%d\n", inm, period, offset)
				results = append(results, struct {
					inm    string
					period int
					offset int
				}{inm, period, offset})
				break
			}
			stateMap[st2] = step
		}
	}

	// Теперь вычислим общий результат
	if len(results) > 0 {
		// Проверим, что у всех offset одинаковы
		baseOffset := results[0].offset
		sameOffset := true
		for _, r := range results {
			if r.offset != baseOffset {
				sameOffset = false
				break
			}
		}
		if !sameOffset {
			fmt.Println("Offsets differ, need CRT. Not implemented here.")
			return
		}
		// Все offset одинаковы.
		allLCM := 1
		for _, r := range results {
			allLCM = lcm(allLCM, r.period)
		}
		n := baseOffset + allLCM
		fmt.Printf("Final Result (minimal n): %d\n", n)
	}
}
