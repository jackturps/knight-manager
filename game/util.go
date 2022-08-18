package game

import (
	"fmt"
	"math/rand"
)

type Number interface {
	int | float64
}

func Max[V Number](item1 V, item2 V) V {
	if item1 > item2 {
		return item1
	} else {
		return item2
	}
}

func Min[V Number](item1 V, item2 V) V {
	if item1 < item2 {
		return item1
	} else {
		return item2
	}
}

func CopySlice[V any](list []V) []V {
	newList := make([]V, len(list))
	copy(newList, list)
	return newList
}

func Exists[V comparable](list []V, item V) bool {
	for _, foundItem := range list {
		if item == foundItem {
			return true
		}
	}
	return false
}

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

func RandomizeOrder[V any] (list []V) []V {
	// Copy to prevent in place modification of input slice.
	listCopy := make([]V, len(list))
	copy(listCopy, list)

	var swapValue V
	for idx, value := range listCopy {
		swapIdx := RandomRange(0, len(listCopy))
		swapValue = listCopy[swapIdx]
		listCopy[swapIdx] = value
		listCopy[idx] = swapValue
	}

	return listCopy
}

var GreenTextCode =     "\x1b[0032m"
var YellowTextCode =    "\x1b[0033m"
var RedTextCode =       "\x1b[0031m"
var DefaultColourCode = "\x1b[0000m"

func ColouredText(colourCode string, text string) string {
	return fmt.Sprintf("%s%s%s", colourCode, text, DefaultColourCode)
}
