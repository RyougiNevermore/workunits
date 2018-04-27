package workunits

import (
	"runtime"
	"time"
)

type reader struct {
	read     *sequence
	written  *sequence
	upstream compositeBarrier
	exec ringExecutor
	ready    bool
}

func newReader(read *sequence, written *sequence, upstream compositeBarrier, exec ringExecutor) *reader {
	return &reader{
		read:     read,
		written:  written,
		upstream: upstream,
		exec: exec,
		ready:    false,
	}
}

func (r *reader) start() {
	r.ready = true
	go r.receive()
}
func (r *reader) stop() {
	r.ready = false
}

func (r *reader) receive() {
	previous := r.read.get()
	idling := 0
	gating := 0
	for {
		lower := previous + 1
		upper := r.upstream.read()
		if lower <= upper {
			r.exec.executor(lower, upper)
			r.read.set(upper)
			previous = upper
		} else if upper = r.written.get(); lower <= upper {
			time.Sleep(time.Microsecond)
			gating++
			idling = 0
		} else if r.ready {
			time.Sleep(time.Millisecond)
			idling++
			gating = 0
		} else {
			break
		}
		runtime.Gosched()
	}
}
