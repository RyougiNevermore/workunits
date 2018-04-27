package workunits

// sequence with cpu mo padding
// only support amd64
type sequence struct {
	value int64
	padding  [7]int64
}

func newSequence() *sequence {
	return &sequence{value: int64(-1)}
}

func (s *sequence) set(value int64) {
	s.value = value
}

func (s *sequence) get() int64 {
	return s.value
}

