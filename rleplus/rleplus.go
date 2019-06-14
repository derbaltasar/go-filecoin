package rleplus

import (
	"sort"
)

func ToBitsetRuns(ints []uint) (firstBit bool, runs []uint) {
	if len(ints) == 0 {
		return
	}

	// Sort our incoming numbers
	sort.Slice(ints, func(i, j int) bool { return ints[i] < ints[j] })

	last := ints[0]

	// Initialize our return value
	firstBit = last == 0
	if !firstBit {
		// first run of zeroes
		runs = append(runs, last)
	}
	runs = append(runs, 1)

	for _, cur := range ints[1:] {
		delta := cur - last
		switch {
		case delta == 1:
			runs[len(runs)-1]++
		case delta > 1:
			// add run of zeroes if there is a gap
			runs = append(runs, cur-last-1)
			runs = append(runs, 1)
		default:
			// repeated number?
		}
		last = cur
	}
	return
}
