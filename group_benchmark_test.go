package workunits

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestDefaultWorkerGroup_Benchmark(t *testing.T) {
	N := 1000000
	var start time.Time
	var dur time.Duration

	start = time.Now()
	benchmarkDefaultGroup(N, 100)
	dur = time.Since(start)
	fmt.Printf("default(100):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkDefaultGroup(N, 10)
	dur = time.Since(start)
	fmt.Printf("default(10):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkDefaultGroup(N, 1)
	dur = time.Since(start)
	fmt.Printf("default(1):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

}

type bunit struct{}

func (u *bunit) Process() {}

func benchmarkDefaultGroup(N int, buffered int64) {

	group := NewDefaultWorkerGroup(buffered, 1024 * 32)
	group.Start()
	for i := 0; i < N; i++ {
		group.Send(&bunit{})
	}
	group.Close()
	group.Sync()
}

func BenchmarkDefaultGroup128(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	benchmarkDefaultGroup(b.N, int64(runtime.NumCPU()))
}

func BenchmarkDefaultGroup16(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	benchmarkDefaultGroup(b.N, 4)
}

func BenchmarkDefaultGroup2(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	benchmarkDefaultGroup(b.N, 1)
}
