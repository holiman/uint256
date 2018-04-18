# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereu/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `fixedbit` being this library. 

As of 2018-03-06:
```
goos: linux
goarch: amd64
pkg: github.com/holiman/uint256

Benchmark_Add/big-2  	100000000	        22.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/fixedbit-2         	300000000	         4.66 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/big-2              	50000000	        23.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/fixedbit-2         	300000000	         4.73 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/fixedbit_of-2      	300000000	         5.12 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/big-2              	10000000	       153 ns/op	     128 B/op	       2 allocs/op
Benchmark_Mul/fixedbit-2         	20000000	        64.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Square/big-2           	10000000	       153 ns/op	     128 B/op	       2 allocs/op
Benchmark_Square/fixedbit-2      	30000000	        49.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/big-2              	100000000	        14.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/fixedbit-2         	1000000000	         1.97 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/big-2               	100000000	        18.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/fixedbit-2          	1000000000	         1.99 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/big-2              	100000000	        19.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/fixedbit-2         	2000000000	         2.01 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/big-2              	200000000	         8.13 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/fixedbit-2         	300000000	         4.44 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/big/n_eq_0-2       	20000000	        93.8 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_192-2     	20000000	       100 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_128-2     	20000000	       105 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_64-2      	20000000	       103 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_0-2       	20000000	        96.9 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/fixedbit/n_eq_0-2  	500000000	         3.97 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/fixedbit/n_gt_192-2         	300000000	         4.59 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/fixedbit/n_gt_128-2         	200000000	         6.53 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/fixedbit/n_gt_64-2          	200000000	         8.41 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/fixedbit/n_gt_0-2           	100000000	        10.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/big/n_eq_0-2                	20000000	        89.5 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_192-2              	20000000	        83.3 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_128-2              	20000000	        84.6 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_64-2               	20000000	        92.4 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_0-2                	20000000	        88.5 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/fixedbit/n_eq_0-2           	500000000	         3.96 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/fixedbit/n_gt_192-2         	300000000	         4.59 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/fixedbit/n_gt_128-2         	200000000	         6.17 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/fixedbit/n_gt_64-2          	200000000	         8.97 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/fixedbit/n_gt_0-2           	100000000	        10.3 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/big-2                 	   50000	     28595 ns/op	   18224 B/op	     191 allocs/op
Benchmark_Exp/large/fixedbit-2            	  100000	     18504 ns/op	      32 B/op	       1 allocs/op
Benchmark_Exp/small/big-2                 	  200000	      7982 ns/op	    7472 B/op	      79 allocs/op
Benchmark_Exp/small/fixedbit-2            	 1000000	      1677 ns/op	      32 B/op	       1 allocs/op
Benchmark_SDiv/large/big-2                	20000000	        94.3 ns/op	      48 B/op	       1 allocs/op
Benchmark_SDiv/large/fixedbit-2           	 2000000	       780 ns/op	     288 B/op	       5 allocs/op
Benchmark_Div/large/big-2                 	 5000000	       340 ns/op	     176 B/op	       3 allocs/op
Benchmark_Div/large/fixedbit-2            	 2000000	       758 ns/op	     288 B/op	       5 allocs/op
Benchmark_Div/small/big-2                 	10000000	       154 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/small/fixedbit-2            	100000000	        16.2 ns/op	       0 B/op	       0 allocs/op


```

The fixed lib wins over big in most cases, with a few exceptions: 

- Division of large numbers. The division algo needs to be replaced with a (pure go) implementation of Knuth's Algorithm D. It
currently wraps `big.Int`

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
