# workunits
Workers pool for Go.

Features:

- Standard and Stable (Default).
- Wild Fly (Coming Soon)
- Wait with Sync.
- Easily Use.

# Getting Started

# Getting Started

### Installing

To start using fastlane, install Go and run `go get`:

```sh
$ go get github.com/pharosnet/workunits
```


### Usage

It's easily to use, default way.

```go

	runtime.GOMAXPROCS(runtime.NumCPU() * 2) 

	group := NewDefaultWorkerGroup(4, 1024 * 32)
	group.Start()
	group.Send(&Unit{...})
	group.Close()
	group.Sync()

```

# TODO

- [] Wild fly mode
- [] Sync() with context.
- [] Work Group Flow

Benchmarks
----------------------------
Any failures cause a panic. Unless otherwise noted, all tests were run using `GOMAXPROCS=runtime.NumCPU() * 2`.

* CPU: `Intel Core i5 @ 2.70 Ghz`
* Operation System: `OS X 10.13.4`
* Go Runtime: `Go 1.10.0`
* Go Architecture: `amd64`

Scenario | Per Operation Time
-------- | ------------------
Default: 1 worker, 1024 * 32 cap, GOMAXPROCS=runtime.NumCPU() * 2| 166 ns/op
Default: 4 worker, 1024 * 32 cap, GOMAXPROCS=runtime.NumCPU() * 2| 194 ns/op
Default: 1 worker, 1024 * 32 cap, GOMAXPROCS=runtime.NumCPU() * 2| 188 ns/op



## Contact

Ryougi Nevermore [@ryougi](https://github.com/RyougiNevermore)

## License

`workunits` source code is available under the MIT [License](/LICENSE).
