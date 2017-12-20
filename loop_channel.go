package workunits

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func NewChanLoop(size int) Loop {
	if size < 1 {
		panic(fmt.Errorf("new chan loop failed, size(%v) is bad", size))
		return nil
	}
	return &chanLoop{
		size:       size,
		unitCh:     nil,
		running:    false,
		statusLock: new(sync.RWMutex),
	}
}

type chanLoop struct {
	size       int
	unitCh     chan Unit
	running    bool
	statusLock *sync.RWMutex
}

func (l *chanLoop) Start() error {
	l.statusLock.Lock()
	defer l.statusLock.Unlock()
	if l.running {
		return loopErrStarted
	}
	if l.unitCh == nil {
		l.unitCh = make(chan Unit, l.size)
	}
	l.running = true
	return nil
}

func (l *chanLoop) Shutdown(ctx context.Context) error {
	l.statusLock.Lock()
	defer l.statusLock.Unlock()
	if !l.running {
		return loopErrStopped
	}
	stopCh := make(chan struct{}, 1)
	stopChClosed := false
	for {
		if !l.running {
			break
		}
		select {
		case <-ctx.Done():
			l.running = false
		case <-stopCh:
			l.running = false
			l.unitCh = nil
		case <-time.After(500 * time.Microsecond):
			if len(l.unitCh) == 0 {
				stopCh <- struct{}{}
				close(stopCh)
				stopChClosed = true
				close(l.unitCh)
			}
		}
	}
	if !stopChClosed {
		close(stopCh)
	}
	return nil
}

func (l *chanLoop) Put(u Unit) (err error) {
	l.statusLock.RLock()
	defer l.statusLock.RUnlock()
	if !l.running {
		return loopErrStopped
	}
	l.unitCh <- u
	return nil
}

func (l *chanLoop) Get() (u Unit, err error) {
	if l.unitCh == nil {
		return nil, loopErrEmpty
	}
	u, ok := <-l.unitCh
	if !ok {
		return nil, loopErrStopped
	}
	if u == nil {
		return nil, loopErrNilWorker
	}
	return u, nil
}

func (l *chanLoop) GetTimeout(timeout time.Duration) (u Unit, err error) {
	if l.unitCh == nil {
		return nil, loopErrEmpty
	}
	select {
	case u, ok := <-l.unitCh:
		if !ok {
			return nil, loopErrStopped
		}
		if u == nil {
			return nil, loopErrNilWorker
		}
		return u, nil
	case <-time.After(timeout):
		return nil, loopErrTimeout
	}
}

func (l *chanLoop) GetDeadline(deadline time.Time) (u Unit, err error) {
	now := time.Now()
	d := deadline.Sub(now)
	if d <= time.Duration(0) {
		return nil, fmt.Errorf("get worker failed, bad deadline(%v), now(%v)", deadline, now)
	}
	return l.GetTimeout(d)
}
