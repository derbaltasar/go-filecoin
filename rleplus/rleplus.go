package rleplus

import (
	"sort"
)

type RLEPlus struct {
	FirstBit bool
	Runs     []uint
}

func NewRLEPlus(ints []uint) (out RLEPlus) {
	if len(ints) == 0 {
		return
	}

	// Sort our incoming numbers
	sort.Slice(ints, func(i, j int) bool { return ints[i] < ints[j] })

	last := ints[0]

	// Initialize our return value
	out.FirstBit = last == 0
	if !out.FirstBit {
		// first run of zeroes
		out.Runs = append(out.Runs, last)
	}
	out.Runs = append(out.Runs, 1)

	for _, cur := range ints[1:] {
		delta := cur - last
		switch {
		case delta == 1:
			out.Runs[len(out.Runs)-1]++
		case delta > 1:
			// add run of zeroes if there is a gap
			out.Runs = append(out.Runs, cur-last-1)
			out.Runs = append(out.Runs, 1)
		default:
			// repeated number?
		}
		last = cur
	}
	return
}
