package workunits

// ring buffer writer.
// only support amd64
type writer struct {
	written  *sequence
	upstream compositeBarrier
	capacity int64
	previous int64
	gate     int64
}

func newWriter(written *sequence, upstream compositeBarrier, capacity int64) *writer {
	if capacity > 0 && (capacity&(capacity-1)) != 0 {
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
		return nil
	}
	return &writer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: int64(-1),
		gate:     int64(-1),
	}
}

func (w *writer) await(next int64) {
	for next-w.capacity > w.gate {
		w.gate = w.upstream.read()
	}
}

func (w *writer) commit(seq int64) {
	w.written.set(seq)
}
