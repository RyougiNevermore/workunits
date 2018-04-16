package workunits

import (
	"fmt"
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
	fmt.Printf("default(100):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkDefaultGroup(N, 1)
	dur = time.Since(start)
	fmt.Printf("default(100):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

}

type bunit struct{}

func (u *bunit) Process() {}

func benchmarkDefaultGroup(N int, buffered int64) {
	group := NewDefaultWorkerGroup(buffered)
	group.Start()
	for i := 0; i < N; i++ {
		group.Send(&bunit{})
	}
	group.Close()
	group.Sync()
}

func BenchmarkDefaultGroup100(b *testing.B) {
	b.ReportAllocs()
	benchmarkDefaultGroup(b.N, 100)
}

func BenchmarkDefaultGroup10(b *testing.B) {
	b.ReportAllocs()
	benchmarkDefaultGroup(b.N, 10)
}

func BenchmarkDefaultGroup1(b *testing.B) {
	b.ReportAllocs()
	benchmarkDefaultGroup(b.N, 1)
}
