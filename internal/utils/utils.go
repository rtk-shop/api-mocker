package utils

import "math/rand/v2"

func RandInt(min, max int) int {
	if min > max {
		min, max = max, min
	}

	return min + rand.IntN(max-min+1)
}

func RandomSliceElement[T any](slice []T) T {
	if len(slice) == 0 {
		panic("slice is empty")
	}
	return slice[rand.IntN(len(slice))]
}
