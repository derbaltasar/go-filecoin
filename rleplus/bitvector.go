package rleplus

const blockSize = 64

type BitVector struct {
	Length uint
	Blocks []uint64

	idx int
}
