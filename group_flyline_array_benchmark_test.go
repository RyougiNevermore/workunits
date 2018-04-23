package workunits

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestArrayBufferedWorkerGroup_Benchmark(t *testing.T) {
	N := 10000
	var start time.Time
	var dur time.Duration

	start = time.Now()
	benchmarkArrayBufferedWorkerGroup(N, 128)
	dur = time.Since(start)
	fmt.Printf("default(100):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkArrayBufferedWorkerGroup(N, 16)
	dur = time.Since(start)
	fmt.Printf("default(10):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkArrayBufferedWorkerGroup(N, 2)
	dur = time.Since(start)
	fmt.Printf("default(1):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

}

type abunit struct{}

func (u *abunit) Process() {}

func benchmarkArrayBufferedWorkerGroup(N int, buffered int64) {

	group := NewArrayBufferedWorkerGroup(buffered)
	group.Start()
	for i := 0; i < N; i++ {
		group.Send(&abunit{})
	}
	group.Close()
	group.Sync()
}

func BenchmarkArrayBufferedWorkerGroup128(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	b.ResetTimer()
	benchmarkArrayBufferedWorkerGroup(b.N, 8)
	b.StopTimer()
}

func BenchmarkArrayBufferedWorkerGroup16(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	b.ResetTimer()
	benchmarkArrayBufferedWorkerGroup(b.N, 16)
}

func BenchmarkArrayBufferedWorkerGroup2(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	b.ResetTimer()
	benchmarkArrayBufferedWorkerGroup(b.N, 2)
}

func BenchmarkArrayBufferedWorkerGroup(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	defer runtime.GOMAXPROCS(1)

	group := NewArrayBufferedWorkerGroup(1024 * 16)
	group.Start()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		group.Send(&abunit{})
	}
	group.Close()
	group.Sync()
	b.StopTimer()
}
