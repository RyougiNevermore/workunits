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

	group := NewDefaultWorkerGroup(4)
	group.Start()
	group.Send(&Unit{...})
	group.Close()
	group.Sync()

```

# TODO

- [] Wild fly mode
- [] Sync() with context.
- [] Work Group Flow

# Performance

## Default Ways

```
$ go test -run none -bench .
``` 

MacBook Pro 13" 2.7 GHz Intel Core i5 (darwin/amd64)

```
10000000	       	363 ns/op	       0 B/op	       0 allocs/op
5000000	       		363 ns/op	       0 B/op	       0 allocs/op
```


## Contact

Ryougi Nevermore [@ryougi](https://github.com/RyougiNevermore)

## License

`workunits` source code is available under the MIT [License](/LICENSE).
