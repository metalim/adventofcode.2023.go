package main

import (
	"fmt"
	"os"
	"sort"
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

type Label int
type Type int

const (
	Labels  = "23456789TJQKA"
	LabelsJ = "J23456789TQKA"

	TypeHigh Type = iota
	TypePair
	TypeTwoPair
	TypeThree
	TypeFullHouse
	TypeFour
	TypeFive
)

type Hand struct {
	cards []Label
	bid   int
	typ   Type
}

func parseHand(line string) *Hand {
	fields := strings.Fields(line)
	hand, bidStr := fields[0], fields[1]

	labels := make([]Label, len(hand))
	for i, card := range hand {
		labels[i] = Label(strings.IndexRune(Labels, card))
	}
	bid, err := strconv.Atoi(bidStr)
	catch(err)

	return &Hand{
		cards: labels,
		bid:   bid,
		typ:   getType(labels),
	}
}

func getType(cards []Label) Type {
	types := make(map[Label]int)
	for _, c := range cards {
		types[c]++
	}
	switch len(types) {
	case 1:
		return TypeFive
	case 2:
		for _, v := range types {
			if v == 4 {
				return TypeFour
			}
		}
		return TypeFullHouse
	case 3:
		for _, v := range types {
			if v == 3 {
				return TypeThree
			}
		}
		return TypeTwoPair
	case 4:
		return TypePair
	case 5:
		return TypeHigh
	}
	panic(fmt.Sprintf("invalid hand: %v", cards))
}

func part1(lines []string) {
	timeStart := time.Now()
	hands := make([]*Hand, len(lines))
	for i, line := range lines {
		hands[i] = parseHand(line)
	}

	sort.Slice(hands, func(i, j int) bool {
		if hands[i].typ != hands[j].typ {
			return hands[i].typ < hands[j].typ
		}
		for i, c := range hands[i].cards {
			if c != hands[j].cards[i] {
				return c < hands[j].cards[i]
			}
		}
		panic("duplicate hands")
	})

	var sum int
	for i, hand := range hands {
		sum += hand.bid * (i + 1)
	}

	fmt.Println("Part 1:", sum, "\tin", time.Since(timeStart))
}

func parseHandJ(line string) *Hand {
	fields := strings.Fields(line)
	hand, bidStr := fields[0], fields[1]

	labels := make([]Label, len(hand))
	for i, card := range hand {
		labels[i] = Label(strings.IndexRune(LabelsJ, card))
	}
	bid, err := strconv.Atoi(bidStr)
	catch(err)

	return &Hand{
		cards: labels,
		bid:   bid,
		typ:   getTypeJ(labels),
	}
}

func getTypeJ(cards []Label) Type {
	types := make(map[Label]int)
	var jokers int
	for _, c := range cards {
		if c == 0 {
			jokers++
			continue
		}
		types[c]++
	}
	switch len(types) {
	case 0, 1:
		return TypeFive
	case 2:
		for _, v := range types {
			if v+jokers == 4 {
				return TypeFour
			}
		}
		return TypeFullHouse
	case 3:
		for _, v := range types {
			if v+jokers == 3 {
				return TypeThree
			}
		}
		return TypeTwoPair
	case 4:
		return TypePair
	case 5:
		return TypeHigh
	}
	panic(fmt.Sprintf("invalid hand: %v", cards))
}

func part2(lines []string) {
	timeStart := time.Now()
	hands := make([]*Hand, len(lines))
	for i, line := range lines {
		hands[i] = parseHandJ(line)
	}

	sort.Slice(hands, func(i, j int) bool {
		if hands[i].typ != hands[j].typ {
			return hands[i].typ < hands[j].typ
		}
		for i, c := range hands[i].cards {
			if c != hands[j].cards[i] {
				return c < hands[j].cards[i]
			}
		}
		panic("duplicate hands")
	})

	var sum int
	for i, hand := range hands {
		sum += hand.bid * (i + 1)
	}

	fmt.Println("Part 2:", sum, "\tin", time.Since(timeStart))
}
