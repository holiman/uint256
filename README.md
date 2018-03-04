# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereu/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `bit` being this library. 

As of 2018-03-04:
```
goos: linux
goarch: amd64
pkg: github.com/holiman/fixed256
Benchmark_Add_Bit-2           	300000000	         4.23 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add_Bit2-2          	300000000	         5.51 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add_Big-2           	50000000	        21.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_SubOverflow_Bit-2   	300000000	         5.57 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub_Bit-2           	300000000	         4.41 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub_Big-2           	100000000	        21.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul_Big-2           	10000000	       145 ns/op	     128 B/op	       2 allocs/op
Benchmark_Mul_Bit-2           	20000000	       104 ns/op	       0 B/op	       0 allocs/op
Benchmark_And_Big-2           	100000000	        13.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_And_Bit-2           	2000000000	         1.90 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or_Big-2            	100000000	        17.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or_Bit-2            	2000000000	         1.90 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor_Big-2           	100000000	        17.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor_Bit-2           	2000000000	         1.89 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp_Big-2           	200000000	         7.72 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp_Bit-2           	1000000000	         2.70 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh_Big-2           	20000000	        92.1 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh_Bit-2           	200000000	         7.00 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh_Big-2           	20000000	        78.6 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh_Bit-2           	200000000	         7.05 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp_Big-2           	   50000	     25713 ns/op	   18224 B/op	     191 allocs/op
Benchmark_Exp_Bit-2           	   50000	     30187 ns/op	      32 B/op	       1 allocs/op
Benchmark_ExpSmall_Big-2      	  200000	      7555 ns/op	    7472 B/op	      79 allocs/op
Benchmark_ExpSmall_Bit-2      	  500000	      3024 ns/op	      32 B/op	       1 allocs/op
Benchmark_DivSmall_Big-2      	10000000	       142 ns/op	     128 B/op	       3 allocs/op
Benchmark_DivSmall_Bit-2      	100000000	        14.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_DivLarge_Big-2      	 5000000	       316 ns/op	     176 B/op	       3 allocs/op
Benchmark_DivLarge_Bit-2      	 1000000	      1275 ns/op	       0 B/op	       0 allocs/op

```

The fixed lib wins over big in most cases, with a few exceptions: 

- Division of large numbers. The division algo needs to be replaced with a (pure go) implementation of Knuth's Algorithm D. 
- Large exponentiation. The underlying `Mul` operation probably needs some additional love and nurturing. 

## Help out

If you're interested in low-level algorithms and/or doing optimizations for shaving off nanoseconds, then this is certainly for you!

Choose an operation, and optimize the s**t out of it!

A few rules, though, to help your PR get approved:

- Do not optimize for 'best-case'/'most common case' at the expense of worst-case. 
- We'll hold off on `asm` for a while, until the algos and interfaces are finished in a first version.

Also, any help in improving the test framework, e.g. by improving the random testing stuff is very highly appreciated. 

