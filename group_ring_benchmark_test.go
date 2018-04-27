package workunits

import (
	"testing"
	"runtime"
)

func BenchmarkRingWorkerGroup128(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	benchmarkRingWorkerGroup(b.N, int64(runtime.NumCPU()))
}

func BenchmarkRingWorkerGroup16(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	benchmarkRingWorkerGroup(b.N, int64(4))
}

func BenchmarkRingWorkerGroup2(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	benchmarkRingWorkerGroup(b.N, int64(1))
}


func benchmarkRingWorkerGroup(N int, workers int64) {
	group := NewRingWorkerGroup(workers, 1024 * 32, 1024 * 32)
	group.Start()
	for i := 0 ; i < N ; i ++ {
		group.Send(&ringUnitb{})
	}
	group.Close()
	group.Sync()
}

type ringUnitb struct {

}

func (u *ringUnitb) Process()  {}