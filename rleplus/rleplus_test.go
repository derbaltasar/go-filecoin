package rleplus_test

import (
	"testing"

	"github.com/filecoin-project/go-filecoin/rleplus"
	tf "github.com/filecoin-project/go-filecoin/testhelpers/testflags"
	"gotest.tools/assert"
)

func TestRleplus(t *testing.T) {
	tf.UnitTest(t)

	t.Run("Generates runs correctly", func(t *testing.T) {
		testCases := []struct {
			ints  []uint
			first bool
			runs  []uint
		}{
			// empty
			{},

			// leading with ones
			{[]uint{0}, true, []uint{1}},
			{[]uint{0, 1}, true, []uint{2}},
			{[]uint{0, 0xffffffff, 0xffffffff + 1}, true, []uint{1, 0xffffffff - 1, 2}},

			// leading with zeroes
			{[]uint{1}, false, []uint{1, 1}},
			{[]uint{2}, false, []uint{2, 1}},
			{[]uint{10, 11, 13, 20}, false, []uint{10, 2, 1, 1, 6, 1}},
			{[]uint{10, 11, 11, 13, 20, 10, 11, 13, 20}, false, []uint{10, 2, 1, 1, 6, 1}},
		}

		for _, testCase := range testCases {
			first, runs := rleplus.ToBitsetRuns(testCase.ints)
			assert.Equal(t, testCase.first, first)
			assert.Equal(t, len(testCase.runs), len(runs))
			for idx, runLength := range testCase.runs {
				assert.Equal(t, runLength, runs[idx])
			}
		}
	})
}
