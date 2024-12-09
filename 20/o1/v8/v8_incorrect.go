/*
```
➜ go run ./o1/v7 input1.txt
Part 2 Answer: no solution
```
что неверно.

Давай пока без суммирования. Агрегатор и подсистемы ты нашёл правильно.
Найди и выведи период цикла для каждой подсистемы отдельно. Не делай допущений
о количестве нажатий, жми кнопку пока не определишь период для каждой подсистемы (2000 — это мало).

Как объединить эти циклы без перебора "от 1 до миллиона" подумаем потом.
Миллиона там недостаточно, т.к. результат — 15-значное число.
Повторю вводные: результат существует, длинное число, помещающееся в int;
допущений о числе нажатий для поиска цикла каждой подсистемы делать нельзя,
т.к. рискуешь преждевременно прекратить поиски

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
Задача: Определить период для каждой подсистемы (каждого входа aggregator) без заранее заданного лимита.

Подход:
1. Ищем модуль-broadcaster и модуль-aggregator (конъюнкцию, ведущую на rx).
2. Каждая подсистема - это путь от одного из выходов broadcaster до одного входа aggregator.
   Чтобы найти период сигнала на входе aggregator, мы будем симулировать нажатия кнопки,
   фиксируя состояние соответствующего входа aggregator после каждого нажатия.
3. Для нахождения цикла без лимита используем алгоритм Брента:
   - Он позволяет найти длину цикла (период) и смещение (offset), не храня всю последовательность.
   - Однако после определения длины цикла lam, нам нужен offset (mu). Для этого перезапустим систему
     и воспользуемся логикой определения начала цикла.
4. Сделаем это для каждого входа aggregator.
5. Выведем периоды и смещения для каждого входа.

Допущения:
- Структура данных корректна: один aggregator, один broadcaster, достижимость есть.
- Результат существует, значит цикл найдётся.
- Число нажатий может быть большим, но предполагаем, что производительность и память достаточно.

Описание алгоритма Брента:
- Инициализируем:
  power = 1, lam = 1
  Берём первое значение x0, tortoise=x0, hare = x0.
  Сдвигаем hare на один шаг (чтобы tortoise=hare были в разных точках).
- Пока tortoise != hare:
  - Если lam == power:
    tortoise = hare
    power = power * 2
    lam = 0
  - Двигаем hare на один шаг вперёд
  - lam++
- Когда нашли совпадение tortoise==hare, lam - длина цикла.

Чтобы найти mu (offset):
- Перезапускаем генерацию.
- Двигаем hare вперёд на lam шагов.
- Тortoise в начале, hare на lam шагах вперёд.
- Пока tortoise != hare:
  - оба двигаются по одному шагу
  - mu++
- mu - смещение начала цикла.

Мы ищем цикл в последовательности состояний конкретного входа aggregator.
На каждом шаге:
  - нажимаем кнопку
  - смотрим aggregator.conInputs[inModule]

Будем реализовывать функцию findCycleForInput для одного входа aggregator.

Замечание:
- Предполагаем, что при очень большом цикле всё равно хватит времени и памяти.
- В реальном случае понадобилась бы оптимизация.

Выводим: входной модуль, период, смещение.
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

// Найдём broadcaster
func findBroadcaster(modules map[string]*Module) string {
	for n, m := range modules {
		if m.mtype == Broadcaster {
			return n
		}
	}
	return ""
}

// Найдём aggregator: conjunction ведущий к rx
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

// Генерирует последовательность для входа inMod: при каждом вызове даёт следующее значение
// Нужно каждый раз стартовать с reset. Но для алгоритма Брента нам нужна последовательность подряд.
// Мы можем хранить состояние modules снаружи. Для поиска цикла нам придётся:
//   - Сделать отдельную функцию, которая каждый вызов даёт следующее значение, нажимая кнопку.
type sequenceGenerator struct {
	modules map[string]*Module
	aggName string
	inMod   string
}

func (g *sequenceGenerator) nextValue() PulseType {
	pressButton(g.modules)
	agg := g.modules[g.aggName]
	return agg.conInputs[g.inMod]
}

// Брент для поиска цикла. Для этого нам нужен только nextValue.
// Но нам нужно будет ещё раз пройтись, чтобы найти mu.
// Brent steps:
// 1 pass: find lambda
// 2 pass: find mu

func brentCycleDetection(next func() PulseType) (int, int) {
	// Псевдокод Брента:
	// power=1; lam=1;
	// tortoise=nextValue()
	// hare=nextValue()
	// while tortoise!=hare:
	//   if power==lam:
	//     tortoise=hare
	//     power*=2
	//     lam=0
	//   hare=nextValue()
	//   lam++
	// нашли цикл длины lam
	// теперь найти mu:
	// reset sequence, проделать то же самое чтобы получить первую лямбда последовательность?
	// Проблема: nextValue меняет состояние.
	// Нам нужен способ заново начать генерацию. Придётся возвращать шаги назад нельзя.
	// Решение: мы не можем reset внутри brent, т.к. у нас только одна попытка.
	// Но пользователь разрешил все средства.
	// Тогда сделаем следующий трюк:
	// Сначала соберём всю последовательность tortoise/hare?
	// Нельзя, sequence может быть огромной.

	// Брент не требует хранить все значения, но для mu нужно сбросить и пройти заново.
	// Мы не можем повторить nextValue с нуля без reset. Решение:
	// Разобьём цикл поиска на две фазы:
	// 1) Определим lam, не сбрасывая.
	// 2) Сохраним lam, затем сделаем reset внешне и вторым вызовом brentPass
	//    чтобы найти mu.
	//
	// Но мы не знаем mu без lam. Mu определяется после нахождения lam.
	//
	// Алгоритм Брента:
	// После цикла у нас есть lam.
	// Чтобы найти mu:
	//   Снова reset.
	//   hare уходит на lam шагов вперёд
	//   tortoise в начале
	//   двигаем по одному пока не совпадут
	// mu - это кол-во шагов до совпадения.

	// Значит, тут мы реализуем только поиск lam.
	// Возвращаем lam, и число шагов сделанных totalSteps.
	// totalSteps нужно чтобы повторить потом процедуру для mu.

	power := 1
	lam := 1
	tortoise := next() // f(1)
	hare := next()     // f(2)
	steps := 2
	for tortoise != hare {
		if power == lam {
			tortoise = hare
			power *= 2
			lam = 0
		}
		hare = next()
		lam++
		steps++
	}
	// lam - длина цикла, steps - сколько шагов сделано
	return lam, steps
}

// Найдём mu:
// Знаем lam.
// Чтобы найти mu:
// reset sequence
// hare проходит lam шагов
// tortoise=f(1), hare=f(1+lam)
// пока tortoise!=hare:
//
//	шаг tortoise, шаг hare
//	mu++
//
// вернуть mu
func findMu(modules map[string]*Module, aggName, inMod string, lam int) int {
	reset(modules)
	g := sequenceGenerator{modules, aggName, inMod}
	// hare на lam шагов вперёд
	for i := 0; i < lam; i++ {
		g.nextValue()
	}
	// теперь tortoise и hare двигаются вместе
	tortoiseVal := sequenceGenerator{modules: modules, aggName: aggName, inMod: inMod}
	reset(modules)
	tortoiseVal.modules = modules
	tortoiseVal.aggName = aggName
	tortoiseVal.inMod = inMod
	// перезапуск tortoise
	reset(modules)
	tortoiseVal.modules = modules
	tortoiseVal.aggName = aggName
	tortoiseVal.inMod = inMod

	tortoiseVal.nextValue()                               // сдвиг на 1, чтобы идти с тех же индексов
	hareVal := sequenceGenerator{modules, aggName, inMod} // уже ушёл на lam шагов
	// НО мы только что reset. Надо понимать, при reset мы потеряли состояние hare.
	// Значит надо заново повторить lam шагов для hare.

	reset(modules)
	hareVal.modules = modules
	hareVal.aggName = aggName
	hareVal.inMod = inMod
	for i := 0; i < lam; i++ {
		hareVal.nextValue()
	}

	mu := 0
	t := tortoiseVal.nextValue() // f(1)
	h := hareVal.nextValue()     // f(1+lam)
	mu++
	for t != h {
		t = tortoiseVal.nextValue()
		h = hareVal.nextValue()
		mu++
	}
	return mu
}

// Упростим процедуру нахождения mu и lam в две фазы:
// 1) Найти lam с одним прогоном Брента.
// 2) Зная lam, найти mu отдельным процедурным прогоном.
//
// Но есть проблема: для brentCycleDetection нужен единичный доступ к next().
// Мы не можем "промотать" назад.
// Решение: Нам нужен метод, который позволяет перезапустить генерацию.
// Применим такой подход:
// - Для lam используем brentCycleDetection один раз на fresh reset.
// - Потом отдельно для mu делаем ещё один reset и ещё одна логика:
//   Сначала lam шагов hare, потом движемся вместе.
//
// Нужно аккуратно: brentCycleDetection нельзя повторно вызвать, он использует nextValue,
// а это изменяет состояние modules. После нахождения lam, система уже далеко.
// Надо сделать так:
//   Сначала находим lam и точку столкновения. Чтобы не портить состояние, сначала находим lam и offset на одном прогоне.
//   Но мы не можем найти offset (mu) без второго прохода.
// Мы и так решили перезапустить систему для нахождения mu, значит состояние после первого прогона неважно.
//
// Итого решение:
// - Сначала один прогон для lam
// - Потом reset и второй прогон для mu (lam известен)
// Но для mu нам нужен доступ к start of cycle. Брент позволяет найти mu без повторного сложного прохода:
// мы просто делаем шаги как описано.

func findCycleForInput(modules map[string]*Module, aggName, inMod string) (int, int) {
	// 1) Найдём lam
	reset(modules)
	g := sequenceGenerator{modules, aggName, inMod}
	lam, _ := brentCycleDetection(g.nextValue)

	// 2) Найдём mu
	mu := findMu(modules, aggName, inMod, lam)

	return lam, mu
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

	// Найдём период для каждого входа aggregator
	for _, inm := range aggMod.conInputList {
		lam, mu := findCycleForInput(modules, aggregator, inm)
		fmt.Printf("Subsystem input %s: period=%d, offset=%d\n", inm, lam, mu)
	}
}
