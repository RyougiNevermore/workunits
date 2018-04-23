package workunits

import (
	"context"
	"fmt"
	"github.com/pharosnet/flyline"
	"os"
	"reflect"
	"sync"
)

// flyline ArrayBuffered worker
type arrayBufferedWorker struct {
	units flyline.Buffer
}

func (w *arrayBufferedWorker) Work() {
	go func(w *arrayBufferedWorker) {
		for {
			iUnit, ok := w.units.Recv()
			if !ok {
				break
			}
			unit, isUnit := iUnit.(Unit)
			if !isUnit {
				panic(fmt.Errorf("arrayBufferedWorker do Work failed. Want Unit{}, but %v recived \n", reflect.TypeOf(iUnit)))
				os.Exit(1)
			}
			unit.Process()
		}
	}(w)
}

// new flyline ArrayBuffered worker group, base on standard channel.
func NewArrayBufferedWorkerGroup(capacity int64) WorkerGroup {
	if capacity > 0 && (capacity&(capacity-1)) != 0 {
		panic("The array capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
		return nil
	}
	group := new(arrayBufferedWorkerGroup)
	group.sts = new(status)
	group.mutex = new(sync.Mutex)
	group.units = flyline.NewArrayBuffer(capacity)
	return group
}

// flyline ArrayBuffered worker group
type arrayBufferedWorkerGroup struct {
	units     flyline.Buffer
	sts       *status
	mutex     *sync.Mutex
	workerNum int64
}

func (g *arrayBufferedWorkerGroup) Start() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isRunning() {
		err = fmt.Errorf("worker group start failed, it is running")
		return
	}
	workerNum := g.workerNum
	for i := int64(0); i < workerNum; i++ {
		worker := new(arrayBufferedWorker)
		worker.units = g.units
		worker.Work()
	}
	g.sts.setRunning()
	return
}

func (g *arrayBufferedWorkerGroup) Close() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isClosed() {
		err = fmt.Errorf("worker group close failed, it is closed")
		return
	}
	g.sts.setClosed()
	g.units.Close()
	return
}

func (g *arrayBufferedWorkerGroup) Sync() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.sts.isRunning() {
		err = fmt.Errorf("worker group sync failed, it is running, please close first")
		return
	}
	g.units.Sync(context.Background())
	return
}

func (g *arrayBufferedWorkerGroup) Send(u Unit) (err error) {
	if g.sts.isClosed() {
		err = fmt.Errorf("worker group send unit to worker failed, it is closed")
		return
	}
	g.units.Send(u)
	return
}
