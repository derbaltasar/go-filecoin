package rleplus_test

import (
	"testing"

	"github.com/filecoin-project/go-filecoin/rleplus"
	tf "github.com/filecoin-project/go-filecoin/testhelpers/testflags"
	"gotest.tools/assert"
)

func TestRuns(t *testing.T) {
	tf.UnitTest(t)

	t.Run("Generates runs correctly", func(t *testing.T) {
		testCases := []struct {
			ints     []uint
			expected rleplus.RLEPlus
		}{
			// empty
			{},

			// leading with ones
			{[]uint{0}, rleplus.RLEPlus{true, []uint{1}}},
			{[]uint{0, 1}, rleplus.RLEPlus{true, []uint{2}}},
			{[]uint{0, 0xffffffff, 0xffffffff + 1}, rleplus.RLEPlus{true, []uint{1, 0xffffffff - 1, 2}}},

			// leading with zeroes
			{[]uint{1}, rleplus.RLEPlus{false, []uint{1, 1}}},
			{[]uint{2}, rleplus.RLEPlus{false, []uint{2, 1}}},
			{[]uint{10, 11, 13, 20}, rleplus.RLEPlus{false, []uint{10, 2, 1, 1, 6, 1}}},
			{[]uint{10, 11, 11, 13, 20, 10, 11, 13, 20}, rleplus.RLEPlus{false, []uint{10, 2, 1, 1, 6, 1}}},
		}

		for _, testCase := range testCases {
			actual := rleplus.NewRLEPlus(testCase.ints)
			assert.Equal(t, testCase.expected.FirstBit, actual.FirstBit)
			assert.Equal(t, len(testCase.expected.Runs), len(actual.Runs))
			for idx, runLength := range testCase.expected.Runs {
				assert.Equal(t, runLength, actual.Runs[idx])
			}
		}
	})
}
