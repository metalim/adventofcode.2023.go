/*
БИНГО!
```
➜ go run ./o1/v10 input1.txt
Subsystem input dh: period=4001, offset=1
Subsystem input mk: period=4091, offset=1
Subsystem input rn: period=3923, offset=1
Subsystem input vf: period=3847, offset=1
```
Теперь добавь в программу код вычисления результата по этим циклам.
Перебирать ничего не требуется, это уже математика.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
Дополняем программу:
1. У нас есть период и смещение для каждого входа подсистемы: period_i и offset_i.
2. Нужно найти момент времени n, при котором все входы будут в состоянии, приводящем к low у агрегатора.
   Из условия: агрегатор выдаёт low, когда все входы high.
   В предыдущих шагах мы узнали, что при достижении цикла состояние модуля повторяется c периодом.
   Если при offset_i начинается период, значит для повторения состояния (которое привело к регистрации цикла)
   мы хотим найти n, удовлетворя n ≡ offset_i (mod period_i).

   Если все offset_i одинаковы (в примере offset=1 у всех), задача упрощается:
   Ищем n такое, что для всех i:
   n ≡ offset (mod period_i).

   Если offset_i у всех одинаков, пусть offset = offset_i для всех i.
   Тогда n-offset делится на каждый period_i.
   Значит n = offset + k * L, где L - НОК всех period_i.
   Минимальное положительное решение: n = offset + LCM(all periods).

   Если offsets не одинаковы, нужно применить китайскую теорему об остатках. Но в примере все offset=1.

3. Вычисляем НОК всех period_i и n = offset + НОК(periods).
4. Выводим это n.

Проверим на примере:
periods: 4001, 4091, 3923, 3847
offsets: все 1
n = 1 + LCM(4001,4091,3923,3847)

Числа большие, но не проблема — мы просто кодируем нахождение НОК через НОД.

Выводим результат.
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
	var periodsOffsets []struct {
		name           string
		period, offset int
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Ждём строки вида: "Subsystem input X: period=Y, offset=Z"
		if strings.HasPrefix(line, "Subsystem input ") {
			// парсим
			// Пример: "Subsystem input dh: period=4001, offset=1"
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				continue
			}
			left := parts[0]  // "Subsystem input dh"
			right := parts[1] // " period=4001, offset=1"

			leftParts := strings.Split(left, " ")
			inm := leftParts[len(leftParts)-1] // последнее слово - имя входа

			right = strings.TrimSpace(right)
			// right = "period=4001, offset=1"
			rparts := strings.Split(right, ",")
			if len(rparts) < 2 {
				continue
			}
			pstr := strings.TrimSpace(rparts[0]) // period=4001
			ostr := strings.TrimSpace(rparts[1]) // offset=1

			pstr = strings.TrimPrefix(pstr, "period=")
			ostr = strings.TrimPrefix(ostr, "offset=")
			var p, o int
			fmt.Sscanf(pstr, "%d", &p)
			fmt.Sscanf(ostr, "%d", &o)

			periodsOffsets = append(periodsOffsets, struct {
				name           string
				period, offset int
			}{inm, p, o})
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(periodsOffsets) == 0 {
		fmt.Println("No subsystems found in input")
		return
	}

	// Проверим, что все offset одинаковы (для упрощения, как в примере)
	baseOffset := periodsOffsets[0].offset
	for _, po := range periodsOffsets {
		if po.offset != baseOffset {
			fmt.Println("Offsets are not the same, need CRT")
			// В общем случае надо решать систему сравнений
			// но в примере все offset=1, так что выходим
			return
		}
	}

	// Вычисляем НОК всех периодов
	allLCM := 1
	for _, po := range periodsOffsets {
		allLCM = lcm(allLCM, po.period)
	}

	n := baseOffset + allLCM
	fmt.Printf("Result (minimal n): %d\n", n)
}
