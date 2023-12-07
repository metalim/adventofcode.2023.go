package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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
	Labels = "23456789TJQKA"

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

func parseHand(hand string) []Label {
	labels := make([]Label, len(hand))
	for i, card := range hand {
		labels[i] = Label(strings.IndexRune(Labels, card))
	}
	return labels
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
	panic("invalid hand")
}

func part1(lines []string) {
	hands := make([]Hand, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		hands[i].cards = parseHand(fields[0])
		bid, err := strconv.Atoi(fields[1])
		catch(err)
		hands[i].bid = bid
		hands[i].typ = getType(hands[i].cards)
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

	fmt.Println("Part 1:", sum)
}

func part2(lines []string) {
}
