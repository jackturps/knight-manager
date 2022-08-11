package game

import (
	"fmt"
	"math/rand"
)

func RemoveItem[V comparable](list []V, item V) []V {
	// Copy to prevent in place modification of input slice. Turns out append modifies!
	listCopy := make([]V, len(list))
	copy(listCopy, list)

	for idx, other := range list {
		if other == item {
			return append(listCopy[:idx], listCopy[idx+1:]...)
		}
	}
	return listCopy
}

func RandomSelect[V any] (values []V) V {
	return values[rand.Intn(len(values))]
}

// RandomRange generates a random number between min and max. This includes min but excludes max.
func RandomRange(min int, max int) int {
	if min > max {
		panic(fmt.Sprintf("min(%d) must be less than or equal to max(%d)", min, max))
	}
	return rand.Intn(max - min) + min
}
