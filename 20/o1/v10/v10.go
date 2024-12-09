/*
подумай ГОЛОВОЙ. ты ищешь циклы у независимых подсистем, сохраняя состояние ВСЕЙ системы.
Это разве адекватно? Сохраняй в хешмапе состояние только конкретной подсистемы.
Найди все узлы, которые принадлежат каждой подсистеме. И сохраняй их состояние
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
Задача:
Определить период для каждой подсистемы (каждого входа aggregator), сохраняя состояние ТОЛЬКО этой подсистемы.

Подсистема:
- Начинается от одного из выходов broadcaster.
- Заканчивается на одном входе aggregator.
- Может быть разветвлённой, но не пересекается с другими подсистемами у aggregator (по условию задачи).

Подход:
1. Найти broadcaster и aggregator.
2. Для каждого входа aggregator:
   - Определить набор модулей, принадлежащих этой подсистеме.
     Для этого:
       * Идём от входного модуля (который подключён к aggregator) назад к broadcaster.
       * Собираем все узлы, которые влияют на этот вход.
         Можно сделать обратный обход по графу, начиная от входного модуля aggregator:
         - Идём по обратным рёбрам пока не достигнем одного из выходов broadcaster.
       * Все модули, достижимые от этого входа (в обратном направлении) до одного из выходов broadcaster — часть подсистемы.

3. Когда набор модулей для подсистемы найден, при симуляции мы будем сохранять состояние только этих модулей:
   - flip-flop (on/off)
   - conjunction (входы)
   Normal модули состояния не хранят.

4. Симулируем нажатия:
   - После каждого нажатия сохраняем текущее состояние подсистемы.
   - Если состояние повторилось, цикл найден:
     period = текущий_step - старый_step_появления
     offset = старый_step_появления

5. Вывести период и offset для каждой подсистемы.

Предполагаем, что каждая подсистема имеет единственный путь к broadcaster (или их можно однозначно выделить).
Если есть несколько путей к broadcaster, все узлы этих путей входят в подсистему.

Пояснение:
- Сначала построим обратный граф.
- Для каждого входа aggregator выполним поиск в обратном графе до узлов, подключенных к broadcaster.
- Узлы, посещённые в этом обратном поиске — это подсистема.
- При симуляции: после нажатия кнопки сохраняем состояние только подсистемных модулей.
- Как только найдём повтор состояния — цикл определён.

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

// Обратный граф
func buildReverseGraph(modules map[string]*Module) map[string][]string {
	rg := make(map[string][]string)
	for src, m := range modules {
		for _, d := range m.dest {
			rg[d] = append(rg[d], src)
		}
	}
	return rg
}

// Найдём узлы подсистемы для конкретного входа aggregator:
//   - startNode = имя входного модуля (из conInputList), этот модуль — вход aggregator
//   - идём по обратному графу до тех модулей, куда ведёт broadcaster (или сам broadcaster)
//
// все достижимые узлы - подсистема
func findSubsystemNodes(modules map[string]*Module, reverseGraph map[string][]string, broadcaster string, startNode string) map[string]bool {
	subNodes := make(map[string]bool)
	// Будем идти вверх до тех пор, пока не достигнем модуля, который является потомком broadcaster
	// broadcaster -> ... -> startNode (в прямом направлении)
	// Нам нужен обратный обход: startNode -> ... -> broadcaster
	// В случае нескольких входов broadcaster надо определить, когда останавливаемся?
	// Если достигли broadcaster — отлично.
	// Если нет, но broadcaster посылает сигнал в некоторые узлы, то достигнем их?
	// В условии сказано, что структура подобна примеру. Предполагаем, что в конце всегда достигнем broadcaster или его детей.
	// На практике: если мы достигли узла broadcaster или узла на который broadcaster непосредственно ссылается — остановимся?

	// Но в примере broadcaster напрямую ведёт в bt, rc, qs, qt — это корни подсистем.
	// Значит если мы достигли узла, который является прямым потомком broadcaster, это корень подсистемы.
	// Проверим это:
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
			// достигли корневой узел подсистемы
			continue
		}
		for _, p := range reverseGraph[cur] {
			stack = append(stack, p)
		}
	}
	return subNodes
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

	// Для каждого входа aggregator найдём подсистему и период
	for _, inm := range aggMod.conInputList {
		// inm — это модуль, который ведёт в aggregator
		// Найдём подсистему
		subNodes := findSubsystemNodes(modules, reverseGraph, broad, inm)

		reset(modules)
		stateMap := make(map[string]int)
		step := 0
		st := subsystemState(modules, subNodes)
		stateMap[st] = step

		for {
			step++
			pressButton(modules)
			// Проверяем состояние подсистемы
			st2 := subsystemState(modules, subNodes)
			if oldStep, ok := stateMap[st2]; ok {
				offset := oldStep
				period := step - oldStep
				fmt.Printf("Subsystem input %s: period=%d, offset=%d\n", inm, period, offset)
				break
			}
			stateMap[st2] = step
		}
	}
}
