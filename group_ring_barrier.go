package workunits

import "math"

type compositeBarrier []*sequence

func newCompositeBarrier(upstream ...*sequence) compositeBarrier {
	if len(upstream) == 0 {
		panic("At least one upstream sequence is required.")
	}
	cursors := make([]*sequence, len(upstream))
	copy(cursors, upstream)
	return compositeBarrier(cursors)
}

func (b compositeBarrier) read() int64 {
	miniSeq := int64(math.MaxInt64)
	for _, seq := range b {
		sequence := seq.get()
		if sequence < miniSeq {
			miniSeq = sequence
		}
	}
	return miniSeq
}
