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
		workerCh:   nil,
		running:    false,
		statusLock: new(sync.RWMutex),
	}
}

type chanLoop struct {
	size       int
	workerCh   chan Worker
	running    bool
	statusLock *sync.RWMutex
}

func (l *chanLoop) Start() error {
	l.statusLock.Lock()
	defer l.statusLock.Unlock()
	if l.running {
		return loopErrStarted
	}
	if l.workerCh == nil {
		l.workerCh = make(chan Worker, l.size)
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
			l.workerCh = nil
		case <-time.After(500 * time.Microsecond):
			if len(l.workerCh) == 0 {
				stopCh <- struct{}{}
				close(stopCh)
				stopChClosed = true
				close(l.workerCh)
			}
		}
	}
	if !stopChClosed {
		close(stopCh)
	}
	return nil
}

func (l *chanLoop) Put(worker Worker) (err error) {
	l.statusLock.RLock()
	defer l.statusLock.RUnlock()
	if !l.running {
		return loopErrStopped
	}
	l.workerCh <- worker
	return nil
}

func (l *chanLoop) Get() (worker Worker, err error) {
	if l.workerCh == nil {
		return nil, loopErrEmpty
	}
	worker, ok := <-l.workerCh
	if !ok {
		return nil, loopErrStopped
	}
	if worker == nil {
		return nil, loopErrNilWorker
	}
	return worker, nil
}

func (l *chanLoop) GetTimeout(timeout time.Duration) (worker Worker, err error) {
	if l.workerCh == nil {
		return nil, loopErrEmpty
	}
	select {
	case worker, ok := <-l.workerCh:
		if !ok {
			return nil, loopErrStopped
		}
		if worker == nil {
			return nil, loopErrNilWorker
		}
		return worker, nil
	case <-time.After(timeout):
		return nil, loopErrTimeout
	}
}

func (l *chanLoop) GetDeadline(deadline time.Time) (worker Worker, err error) {
	now := time.Now()
	d := deadline.Sub(now)
	if d <= time.Duration(0) {
		return nil, fmt.Errorf("get worker failed, bad deadline(%v), now(%v)", deadline, now)
	}
	return l.GetTimeout(d)
}
