package util

import (
	"fmt"
	"math/rand"
	"sort"
)

func GetRandomSet(max, count int) []int {
	var values []int

	if count > max {
		panic(fmt.Sprintf("Invalid args: %v > %v", count, max))
	}

	for i := 0; i < count; i++ {
		v := rand.Intn(max)
		for _, vx := range values {
			if vx <= v {
				v++
			}
		}
		values = append(values, v)
		sort.Ints(values)
		max--
	}
	return values
}
