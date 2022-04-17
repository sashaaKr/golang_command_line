## Getting started

To extract benchmark data:

```sh
tar -xzvf colStatsBenchmarkData.tar.gz -C testdata/
```

to run benchamrk tests:

```sh
go test -bench . -run ^$
go test -bench . -benchtime=10x -run ^$

# profiling cpu
go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu00.pprof
go tool pprof cpu00.ppfor
# in interactive terminal run
top
top -cum
list csv2float
web
quit

# memory profiling
go test -bench . -benchtime=10x -run ^$ -memprofile mem00.pprof
go tool pprof mem00.pprof
top -cum
quit

go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults00m.txt
```

to compare benchmarks

```sh
go get -u -v golang.org/x/perf/cmd/benchstat
```

run tracer tests:

```sh
go test -bench . -benchtime=10x -run ^$ -trace trace01.out
go tool trace trace01.out
```
