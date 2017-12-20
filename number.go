package workunits

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

const number_retry_times = 64
const number_min_value = int64(-1)
const number_max_value = math.MaxInt64

type number struct {
	v          int64
	retryTimes int
	min        int64
	max        int64
	lock       *sync.Mutex
}

func (n *number) Load() int64 {
	return atomic.LoadInt64(&n.v)
}

// add number by cas and mutex when over retry times
func (n *number) Add(i int64) (int64, error) {
	retry := n.retryTimes
	if n.v >= n.max {
		return n.v, fmt.Errorf("add failed, number(%v) is max(%v), ", n.v, n.max)
	}
	for !atomic.CompareAndSwapInt64(&n.v, n.v, n.v+i) {
		retry--
		if retry > 0 {
			continue
		}
		n.lock.Lock()
		n.v = n.v + i
		n.lock.Unlock()
		break
	}
	return n.v, nil
}

// increase number by cas and mutex when over retry times
func (n *number) Inc() (int64, error) {
	return n.Add(int64(1))
}

// set min into number
func (n *number) Reset() {
	n.lock.Lock()
	n.v = n.min
	n.lock.Unlock()
}
