package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	bs, err := os.ReadFile(os.Args[1])
	catch(err)
	lines := strings.Split(string(bs), "\n")

	fmt.Println(lines)
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
