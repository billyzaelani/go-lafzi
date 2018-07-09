// Package sequence provide functionality around finding subsequent
// in streaming data.
package sequence

import (
	"sort"
)

// LIS finds longest increasing subsequent (LIS) in s.
func LIS(s []int) []int {
	var l, k int

	A := make([]int, 2, len(s)+2)
	A[0] = -1000000
	A[1] = 1000000

	seq := make([][]int, 1, len(s))
	seq[0] = []int{}

	for _, x := range s {
		l = sort.Search(len(A), func(i int) bool { return A[i] >= x })
		if A[l] != x {
			l--
		}

		A[l+1] = x

		t := append(seq[l], x)
		if isSet(seq, l+1) {
			seq[l+1] = t
		} else {
			seq = append(seq, t)
		}

		if l+1 > k {
			k++
			A = append(A, 1000000)
			seq = append(seq, []int{})
		}
	}

	return seq[k]
}

// isSet for sequence.
func isSet(arr [][]int, i int) bool {
	return len(arr) > i
}
