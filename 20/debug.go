package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func exportDot(nodes map[string]*Node, fname string) {
	f, err := os.Create(fname)
	catch(err)
	defer f.Close()

	var names [][2]string
	for name := range nodes {
		names = append(names, [2]string{nodes[name].Type + name, name})
	}
	slices.SortFunc(names, func(a, b [2]string) int {
		return strings.Compare(b[0], a[0])
	})
	names = append(names[1:], names[0])
	slices.Reverse(names)

	fmt.Fprintln(f, "digraph aoc {")
	for _, name := range names {
		if len(nodes[name[1]].Dest) == 0 {
			continue
		}
		var dest strings.Builder
		for i, d := range nodes[name[1]].Dest {
			if i > 0 {
				dest.WriteString(", ")
			}
			dest.WriteString("\"")
			dest.WriteString(d)
			if n, ok := nodes[d]; ok {
				if n.Type != Broadcaster {
					dest.WriteString(n.Type)
				}
			} else {
				dest.WriteString("?")
			}
			dest.WriteString("\"")
		}
		fmt.Fprintf(f, "\t\"%s%s\" -> %s;\n", name[1], nodes[name[1]].Type, dest.String())
	}
	fmt.Fprintln(f, "}")
}
