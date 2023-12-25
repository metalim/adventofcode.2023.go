package main

import (
	"fmt"
	"os"
	"regexp"
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
}

type Graph struct {
	Edges map[string]map[string]int
}

var reEdge = regexp.MustCompile(`(\w+)`)

func parse(lines []string) Graph {
	graph := Graph{
		Edges: map[string]map[string]int{},
	}

	for _, line := range lines {
		m := reEdge.FindAllString(line, -1)
		for _, name := range m[1:] {
			if _, ok := graph.Edges[m[0]]; !ok {
				graph.Edges[m[0]] = map[string]int{}
			}
			graph.Edges[m[0]][name] = 1
			if _, ok := graph.Edges[name]; !ok {
				graph.Edges[name] = map[string]int{}
			}
			graph.Edges[name][m[0]] = 1
		}
	}

	// sort.Strings(graph.Nodes)
	// fmt.Println(graph)
	return graph
}

func part1(lines []string) {
	// timeStart := time.Now()
	graph := parse(lines)
	var grouped string
	for node := range graph.Edges {
		grouped = node
		break
	}
	fmt.Println(grouped)
	fmt.Println(len(graph.Edges))
	for len(graph.Edges) > 1 {
		node, n := graph.FindMaxConnectedTo(grouped)
		fmt.Printf("%s: %d, ", node, n)
		graph.Merge(grouped, node)
		// fmt.Println(len(graph.Edges))
	}
	// fmt.Printf("Part 1: \tin %v\n", time.Since(timeStart))
}

func (g Graph) FindMaxConnectedTo(node1 string) (string, int) {
	var maxWeight int
	var maxNode string
	for node2, weight := range g.Edges[node1] {
		if node2 == node1 {
			// panic("circular dependency")
			continue
		}
		if maxWeight < weight {
			maxWeight = weight
			maxNode = node2
		}
	}
	return maxNode, maxWeight
}

func (g Graph) Merge(node1, node2 string) {
	for node3, weight := range g.Edges[node2] {
		// if node3 == node1 {
		// 	continue
		// }
		g.Edges[node1][node3] += weight
		g.Edges[node3][node1] += weight
	}
	delete(g.Edges, node2)
	for _, edges := range g.Edges {
		delete(edges, node2)
	}
}
