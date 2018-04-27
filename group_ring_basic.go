package workunits

import (
	"sync"
	"fmt"
)

// mark(reader and writer) and units , wg
type ringExecutor interface {
	start()
	stop() // with sync
	send(uint Unit)
	executor(lower int64, upper int64)
}

// boss ring
type bossRing struct {
	units []Unit

	mask int64
	nextSeq int64
	wg *sync.WaitGroup
	writer *writer
	reader *reader

	workNum int64
	workers []*workRing
	workNext int64
}

func newBossRing(bossCapacity int64, workCapacity int64, workers int64) *bossRing  {
	if workers <= 0 {
		panic(fmt.Errorf("new boss ring failed, workers is letter than 1"))
		return nil
	}
	boss := new(bossRing)
	boss.units = make([]Unit, bossCapacity)
	boss.mask = bossCapacity - 1
	boss.nextSeq = int64(-1)
	boss.wg = new(sync.WaitGroup)
	// writer and reader
	writerSeq := newSequence()
	writerBarrier := newCompositeBarrier(writerSeq)
	readerSeq := newSequence()
	readerBarrier := newCompositeBarrier(readerSeq)
	boss.reader = newReader(readerSeq, writerSeq, writerBarrier, boss)
	boss.writer = newWriter(writerSeq, readerBarrier, bossCapacity)

	boss.workNum = workers
	boss.workNext = int64(-1)
	boss.workers = make([]*workRing, 0, workers)
	for i := int64(0) ; i < workers ; i ++ {
		boss.workers = append(boss.workers, newWorkRing(workCapacity))
	}
	return boss
}

func (r *bossRing) start()  {
	r.reader.start()
	for _, worker := range r.workers {
		worker.start()
	}
}

func (r *bossRing) stop()  {
	r.reader.stop()
	r.wg.Wait()
	for _, worker := range r.workers {
		worker.stop()
	}
}

func (r *bossRing) send(unit Unit)  {
	r.wg.Add(1)
	next := r.nextSeq + 1
	r.writer.await(next)
	r.units[next&r.mask] = unit
	r.writer.commit(next)
	r.nextSeq = next
}

func (r *bossRing) executor(lower int64, upper int64) {
	for lower <= upper {
		unit := r.units[lower&r.mask]
		r.workNext ++
		r.workers[r.workNext % r.workNum].send(unit)
		r.wg.Done()
		lower ++
	}
}

// work ring
type workRing struct {
	units []Unit
	mask int64
	nextSeq int64
	wg *sync.WaitGroup
	writer *writer
	reader *reader
}

func newWorkRing(capacity int64) *workRing {
	worker := new(workRing)
	worker.units = make([]Unit, capacity)
	worker.mask = capacity - 1
	worker.nextSeq = int64(-1)
	worker.wg = new(sync.WaitGroup)
	// writer and reader
	writerSeq := newSequence()
	writerBarrier := newCompositeBarrier(writerSeq)
	readerSeq := newSequence()
	readerBarrier := newCompositeBarrier(readerSeq)
	worker.reader = newReader(readerSeq, writerSeq, writerBarrier, worker)
	worker.writer = newWriter(writerSeq, readerBarrier, capacity)
	return worker
}


func (r *workRing) start()  {
	r.reader.start()
}


func (r *workRing) stop()  {
	r.reader.stop()
	r.wg.Wait()
}

func (r *workRing) executor(lower int64, upper int64) {
	for lower <= upper {
		unit := r.units[lower&r.mask]
		unit.Process()
		r.wg.Done()
		lower ++
	}
}


func (r *workRing) send(unit Unit)  {
	r.wg.Add(1)
	next := r.nextSeq + 1
	r.writer.await(next)
	r.units[next&r.mask] = unit
	r.writer.commit(next)
	r.nextSeq = next
}