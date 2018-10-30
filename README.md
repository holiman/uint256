# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereu/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `fixedbit` being this library. 

As of 2018-10-30:
```
Benchmark_Add/big-6                 100000000	        21.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/uint256-6         	300000000	         4.47 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/big-6             	100000000	        21.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/uint256-6         	300000000	         4.27 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/uint256_of-6      	300000000	         4.74 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/big-6             	10000000	       146 ns/op	     128 B/op	       2 allocs/op
Benchmark_Mul/uint256-6         	30000000	        58.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Square/big-6          	10000000	       143 ns/op	     128 B/op	       2 allocs/op
Benchmark_Square/uint256-6      	30000000	        46.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/big-6             	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/uint256-6         	2000000000	         1.91 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/big-6              	100000000	        17.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/uint256-6          	2000000000	         1.89 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/big-6             	100000000	        17.3 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/uint256-6         	2000000000	         1.90 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/big-6             	200000000	         7.87 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/uint256-6         	300000000	         3.91 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/big/n_eq_0-6      	20000000	        92.0 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_192-6    	20000000	       103 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_128-6    	20000000	        91.6 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_64-6     	20000000	        90.1 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_0-6      	20000000	        87.5 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/uint256/n_eq_0-6  	500000000	         3.79 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_192-6         	300000000	         4.60 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_128-6         	200000000	         5.80 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_64-6          	200000000	         7.67 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_0-6           	200000000	         9.94 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/big/n_eq_0-6               	20000000	        82.1 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_192-6             	20000000	        77.3 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_128-6             	20000000	        79.7 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_64-6              	20000000	        80.9 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_0-6               	20000000	        81.9 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/uint256/n_eq_0-6           	500000000	         3.85 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_192-6         	300000000	         4.32 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_128-6         	200000000	         5.80 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_64-6          	200000000	         8.58 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_0-6           	100000000	        10.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/big-6                	   50000	     27726 ns/op	   18224 B/op	     191 allocs/op
Benchmark_Exp/large/uint256-6            	  100000	     16855 ns/op	      32 B/op	       1 allocs/op
Benchmark_Exp/small/big-6                	  200000	      7647 ns/op	    7472 B/op	      79 allocs/op
Benchmark_Exp/small/uint256-6            	 1000000	      1606 ns/op	      32 B/op	       1 allocs/op
Benchmark_SDiv/large/big-6               	20000000	        98.8 ns/op	      48 B/op	       1 allocs/op
Benchmark_SDiv/large/uint256-6           	 5000000	       261 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/large/big-6                	 5000000	       318 ns/op	     176 B/op	       3 allocs/op
Benchmark_Div/large/uint256-6            	 5000000	       251 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/small/big-6                	10000000	       142 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/small/uint256-6            	100000000	        14.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mulmod/large/big-6             	 5000000	       320 ns/op	     176 B/op	       3 allocs/op
Benchmark_Mulmod/large/uint256-6         	 1000000	      1272 ns/op	     608 B/op	      11 allocs/op
Benchmark_Mod/large/big-6         	        10000000	       143 ns/op	      48 B/op	       1 allocs/op
Benchmark_Mod/large/uint256-6     	         2000000	       686 ns/op	     352 B/op	       7 allocs/op
Benchmark_Mod/small/big-6         	        20000000	        63.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Mod/small/uint256-6     	       100000000	        16.1 ns/op	       0 B/op	       0 allocs/op
PASS

```

The fixed lib wins over big in most cases, with a few exceptions: 

- Signed division is slower on `uint256`. 
- `MulMod` is slower on `uint256`. 
- `Mod` on large numbers is slower on `uint256`. 

Both `MulMod` and `Mod` currently wraps `big.Int`, which is suboptimal. 

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
