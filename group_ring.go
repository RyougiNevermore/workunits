package workunits

import (
	"sync"
	"fmt"
)

func NewRingWorkerGroup(workerNum int64, bossCapacity int64, workCapacity int64) WorkerGroup {
	if bossCapacity <= int64(0) && workCapacity <= int64(0) && workerNum <= int64(0) {
		panic(fmt.Errorf("new group failed, cap must be bigger than 0"))
	}
	group := new(ringWorkerGroup)
	group.wg = new(sync.WaitGroup)
	group.sts = new(status)
	group.mutex = new(sync.Mutex)
	group.ring = newBossRing(bossCapacity, workCapacity, workerNum)
	return group
}

// ring buffered worker group.
// only support amd64.
// only support no-contended upstream
type ringWorkerGroup struct {
	ring *bossRing
	wg    *sync.WaitGroup
	sts   *status
	mutex *sync.Mutex
}

func (g *ringWorkerGroup) Start() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isRunning() {
		err = fmt.Errorf("worker group start failed, it is running")
		return
	}
	g.ring.start()
	g.sts.setRunning()
	return
}

func (g *ringWorkerGroup) Close() (err error)  {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isClosed() {
		err = fmt.Errorf("worker group close failed, it is closed")
		return
	}
	g.sts.setClosed()
	g.ring.stop()
	return
}

// do nothing, sync() in close()
func (g *ringWorkerGroup) Sync() (err error) {

	return
}

func (g *ringWorkerGroup) Send(u Unit) (err error)  {
	if g.sts.isClosed() {
		err = fmt.Errorf("worker group send unit to worker failed, it is closed")
		return
	}
	g.ring.send(u)
	return
}

