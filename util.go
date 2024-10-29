package matchqueue

import (
	"cmp"
	"sort"
)

// clamp adjusts v in the range of [lower, upper].
func clamp[T cmp.Ordered](v, lower, upper T) T {
	switch {
	case v < lower:
		return lower
	case v > upper:
		return upper
	}
	return v
}

// insertSortedSlice inserts an element into the slice, which is sorted by the given function.
// It inserts the element at the smallest index i in [0, n) at which f(i) is true.
// If f(i) is not satisfied, the element is appended at the last of the slice.
func insertSortedSlice[T any](s []T, e T, f func(int) bool) []T {
	i := sort.Search(len(s), f)

	var z T
	s = append(s, z)
	copy(s[i+1:], s[i:])
	s[i] = e
	return s
}