/*
1. запятые в Go нельзя переносить на новую строку, даже если перед ними многострочный комментарий.
2. в Go нельзя оставлять неиспользуемые переменные. Тут это processPulses, out, sendPulse, startPresses.
3. выводи результат и время по мере заканчивания частей, а не жди финала.
4. и самое главное: брутфорс второй части будет выполняться до окончания жизни вселенной.
Нужно проанализировать последние узлы (возможно граф делится на несколько независимых частей,
объединённых в конце), и на этой базе построить код находящий решение.

Вот пример реальных входных данных:
```
%lh -> mj
%nd -> qf
&pn -> dh, dk, bg, qs, rp, bk, gs
%bk -> rs
%nh -> lh
%hc -> jg, ks
%pt -> gv, jg
&dh -> jz
%jq -> nd
%gv -> jg, mr
%gm -> jv
&zt -> jq, rn, nd, bt, jh, gm
%mz -> dc, zt
%nf -> dm, pn
%bg -> bk
%qt -> qx, xk
%dc -> zt, db
%rc -> gz, jg
%kx -> pn, gj
%mj -> zm
%rs -> pn, dk
%lv -> tb, jg
&mk -> jz
%bt -> pv, zt
%cg -> mz, zt
%pk -> qx
%jd -> lv, jg
%jv -> jh, zt
%ks -> jg, jd
%gs -> bg
broadcaster -> bt, rc, qs, qt
%dm -> rm, pn
%pv -> jq, zt
%db -> zt
%dv -> sl, qx
%qs -> rp, pn
%sr -> hf
%qf -> gm, zt
&jz -> rx
&vf -> jz
%gz -> vj, jg
%mr -> jg
%dk -> kx
&jg -> rc, mk, vj
%qh -> hc, jg
%vj -> qh
%tb -> pt, jg
%rm -> pn
%gj -> pn, nf
%rp -> gs
%xk -> td, qx
%hf -> nh
&rn -> jz
&qx -> lh, vf, hf, nh, sr, mj, qt
%td -> sr, qx
%sl -> pk, qx
%jh -> cg
%zm -> dv, qx
```
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

// Обработка входящего импульса для модуля
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

// Состояние системы: состояние всех флип-флопов и conjunction-модулей
type State struct {
	flipFlops    map[string]bool
	conjunctions map[string][]PulseType // последовательность соответствует conInputList
}

// Превращение текущего состояния модулей в State
func getCurrentState(modules map[string]*Module) State {
	ff := make(map[string]bool)
	conj := make(map[string][]PulseType)
	for name, m := range modules {
		if m.mtype == FlipFlop {
			ff[name] = m.flipOn
		}
		if m.mtype == Conjunction {
			arr := make([]PulseType, len(m.conInputList))
			for i, inp := range m.conInputList {
				arr[i] = m.conInputs[inp]
			}
			conj[name] = arr
		}
	}
	return State{flipFlops: ff, conjunctions: conj}
}

// Применение состояния к модулям
func applyState(modules map[string]*Module, s State) {
	for name, m := range modules {
		if m.mtype == FlipFlop {
			m.flipOn = s.flipFlops[name]
		}
		if m.mtype == Conjunction {
			for i, inp := range m.conInputList {
				m.conInputs[inp] = s.conjunctions[name][i]
			}
		}
	}
}

// Сериализация состояния для хранения в map
func serializeState(s State) string {
	// Сериализуем состояние флип-флопов по имени в лексикографическом порядке
	// и conjunction по имени, затем их входы
	// Это нужно для BFS
	var sb strings.Builder
	sb.WriteString("FF:")
	ffnames := make([]string, 0, len(s.flipFlops))
	for n := range s.flipFlops {
		ffnames = append(ffnames, n)
	}
	var sortStrings func([]string)
	sortStrings = func(a []string) {
		// простой bubble sort достаточно
		for i := 0; i < len(a); i++ {
			for j := i + 1; j < len(a); j++ {
				if a[j] < a[i] {
					a[i], a[j] = a[j], a[i]
				}
			}
		}
	}
	sortStrings(ffnames)
	for _, n := range ffnames {
		if s.flipFlops[n] {
			sb.WriteString(n + "=1,")
		} else {
			sb.WriteString(n + "=0,")
		}
	}

	sb.WriteString(";CONJ:")
	conjnames := make([]string, 0, len(s.conjunctions))
	for n := range s.conjunctions {
		conjnames = append(conjnames, n)
	}
	sortStrings(conjnames)
	for _, n := range conjnames {
		sb.WriteString(n + "=")
		for _, v := range s.conjunctions[n] {
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

// Полная обработка импульсов до устойчивого состояния
func runPulses(modules map[string]*Module, pulses []struct {
	from, to string
	ptype    PulseType
}, countPulses *struct{ low, high int }, stopAtRxLow bool) bool {
	queue := pulses
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		if p.ptype == Low {
			countPulses.low++
		} else {
			countPulses.high++
		}

		if p.to == "rx" && p.ptype == Low && stopAtRxLow {
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

// Нажатие кнопки
func pressButton(modules map[string]*Module, countPulses *struct{ low, high int }, stopAtRxLow bool) bool {
	return runPulses(modules, []struct {
		from, to string
		ptype    PulseType
	}{{from: "button", to: "broadcaster", ptype: Low}}, countPulses, stopAtRxLow)
}

// Сброс в начальное состояние
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

		modules[name] = &Module{
			mtype:     mtype,
			name:      name,
			dest:      dests,
			conInputs: make(map[string]PulseType),
		}
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

	// Часть 1
	startPart1 := time.Now()
	reset(modules)
	var countPulses struct{ low, high int }
	for i := 0; i < 1000; i++ {
		pressButton(modules, &countPulses, false)
	}
	part1Low, part1High := countPulses.low, countPulses.high
	part1Answer := int64(part1Low) * int64(part1High)
	fmt.Println("Part 1 Answer:", part1Answer)
	fmt.Println("Time Part 1:", time.Since(startPart1))

	// Часть 2
	// Нужно минимальное число нажатий на кнопку, чтобы rx получил low.
	// Реализуем BFS по состояниям. Каждое нажатие кнопки – переход в новое состояние.
	// При обработке импульсов проверяем, отправился ли low на rx.

	startPart2 := time.Now()
	reset(modules)
	visited := make(map[string]bool)
	type queueItem struct {
		state State
		steps int
	}
	initialState := getCurrentState(modules)
	visited[SerializeState(initialState)] = true
	q := []queueItem{{state: initialState, steps: 0}}
	foundSteps := -1

	// При каждом нажатии: полностью отрабатываем импульсы. Если low->rx найден, возвращаем steps+1
	// Иначе получаем новый стейt, если не видели его – в очередь.

	for len(q) > 0 && foundSteps < 0 {
		cur := q[0]
		q = q[1:]
		applyState(modules, cur.state)
		var cp struct{ low, high int }
		if pressButton(modules, &cp, true) {
			foundSteps = cur.steps + 1
			break
		}
		// после полной обработки импульсов новый стабильный стейт
		ns := getCurrentState(modules)
		key := SerializeState(ns)
		if !visited[key] {
			visited[key] = true
			q = append(q, queueItem{state: ns, steps: cur.steps + 1})
		}
	}

	fmt.Println("Part 2 Answer:", foundSteps)
	fmt.Println("Time Part 2:", time.Since(startPart2))
}

// SerializeState вынесена, чтобы не плодить код
func SerializeState(s State) string {
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
	conjnames := make([]string, 0, len(s.conjunctions))
	for n := range s.conjunctions {
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
		for _, v := range s.conjunctions[n] {
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
