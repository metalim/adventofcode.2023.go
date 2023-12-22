package main

import "strconv"

func toInts(ss []string) []int {
	is := make([]int, len(ss))
	for i, s := range ss {
		n, err := strconv.Atoi(s)
		catch(err)
		is[i] = n
	}
	return is
}
