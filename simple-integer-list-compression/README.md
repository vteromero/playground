Prototype to benchmark a simple integer compression algorithm. [Here](https://vteromero.github.io/benchmarking-simple-integer-list-compression/) is the corresponding blog post.

### Installation

```
go get github.com/dataence/encoding
go get -t github.com/vteromero/playground
```

### How to run the benchmarks

```
go test -bench=.
```

### How to use `compare-compression-ratio`

Firstly, you need to build the binary. Just type the following:

```
cd cmd/compare-compression-ratio
go build .
```

Once `compare-compression-ratio` binary is generated, you can run it like this:

```
./compare-compression-ratio -sizes=1000,100000
```

You can check all the available options with the `help` flag:

```
./compare-compression-ratio -help
```
