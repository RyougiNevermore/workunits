package workunits

import (
	"fmt"
	"sync"
)

// standard channel worker
type defaultWorker struct {
	units chan Unit
	wg    *sync.WaitGroup
}

func (w *defaultWorker) Work() {
	go func(w *defaultWorker) {
		for {
			unit, ok := <-w.units
			if !ok {
				break
			}
			unit.Process()
			w.wg.Done()
		}
	}(w)
}

// new default worker group, base on standard channel.
func NewDefaultWorkerGroup(workerNum int64, capacity int64) WorkerGroup {
	if capacity <= int64(0) && workerNum <= int64(0) {
		panic(fmt.Errorf("new group failed, cap must be bigger than 0"))
	}
	group := new(defaultWorkerGroup)
	group.workerNum = workerNum
	group.wg = new(sync.WaitGroup)
	group.sts = new(status)
	group.mutex = new(sync.Mutex)
	group.units = make(chan Unit, capacity)
	return group
}

// standard channel worker group
type defaultWorkerGroup struct {
	workerNum int64
	units chan Unit
	wg    *sync.WaitGroup
	sts   *status
	mutex *sync.Mutex
}

func (g *defaultWorkerGroup) Start() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isRunning() {
		err = fmt.Errorf("worker group start failed, it is running")
		return
	}
	for i := int64(0); i < g.workerNum; i++ {
		worker := new(defaultWorker)
		worker.units = g.units
		worker.wg = g.wg
		worker.Work()
	}
	g.sts.setRunning()
	return
}

func (g *defaultWorkerGroup) Close() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isClosed() {
		err = fmt.Errorf("worker group close failed, it is closed")
		return
	}
	g.sts.setClosed()
	close(g.units)
	return
}

func (g *defaultWorkerGroup) Sync() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isRunning() {
		err = fmt.Errorf("worker group sync failed, it is running, please close first")
		return
	}
	g.wg.Wait()
	return
}

func (g *defaultWorkerGroup) Send(u Unit) (err error) {
	if g.sts.isClosed() {
		err = fmt.Errorf("worker group send unit to worker failed, it is closed")
		return
	}
	g.wg.Add(1)
	g.units <- u
	return
}
