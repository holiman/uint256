# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereu/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `fixedbit` being this library. 

As of 2018-03-05:
```
goos: linux
goarch: amd64
pkg: github.com/holiman/fixed256
Benchmark_Add/big-2         	50000000	        22.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/fixedbit-2    	300000000	         4.28 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/big-2         	100000000	        22.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/fixedbit-2    	300000000	         4.45 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/fixedbit_of-2 	300000000	         5.39 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/big-2         	10000000	       146 ns/op	     128 B/op	       2 allocs/op
Benchmark_Mul/fixedbit-2    	20000000	        61.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/big-2         	100000000	        14.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/fixedbit-2    	1000000000	         1.98 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/big-2          	100000000	        18.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/fixedbit-2     	2000000000	         1.95 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/big-2         	100000000	        17.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/fixedbit-2    	2000000000	         1.93 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/big-2         	200000000	         8.08 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/fixedbit-2    	500000000	         3.00 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/big-2         	20000000	        93.1 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/fixedbit-2    	500000000	         3.96 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/big-2         	20000000	        80.1 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/fixedbit-2    	300000000	         4.19 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/big-2   	   50000	     25330 ns/op	   18224 B/op	     191 allocs/op
Benchmark_Exp/large/fixedbit-2         	  100000	     20600 ns/op	      32 B/op	       1 allocs/op
Benchmark_Exp/small/big-2              	  200000	      7695 ns/op	    7472 B/op	      79 allocs/op
Benchmark_Exp/small/fixedbit-2         	 1000000	      1979 ns/op	      32 B/op	       1 allocs/op
Benchmark_Div/large/big-2              	 5000000	       308 ns/op	     176 B/op	       3 allocs/op
Benchmark_Div/large/fixedbit-2         	 1000000	      1251 ns/op	       0 B/op	       0 allocs/op
Benchmark_Div/small/big-2              	10000000	       143 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/small/fixedbit-2         	100000000	        14.7 ns/op	       0 B/op	       0 allocs/op


```

The fixed lib wins over big in most cases, with a few exceptions: 

- Division of large numbers. The division algo needs to be replaced with a (pure go) implementation of Knuth's Algorithm D. 

## Help out

If you're interested in low-level algorithms and/or doing optimizations for shaving off nanoseconds, then this is certainly for you!

Choose an operation, and optimize the s**t out of it!

A few rules, though, to help your PR get approved:

- Do not optimize for 'best-case'/'most common case' at the expense of worst-case. 
- We'll hold off on go assembly for a while, until the algos and interfaces are finished in a 'good enough' first version. After that, it's assembly time. 

Also, any help in improving the test framework, e.g. by improving the random testing stuff is very highly appreciated. 

### Doing benchmarks

To do a simple benchmark for everything, do

```
go test -run - -bench . -benchmem

```

To see the difference between a branch and master, for a particular benchmark, do

```
git checkout master
go test -run - -bench Benchmark_Lsh -benchmem -count=10 > old.txt

git checkout opt_branch
go test -run - -bench Benchmark_Lsh -benchmem -count=10 > new.txt

benchstat old.txt new.txt

```
