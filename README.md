# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereum/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `uint256` being this library. 

As of 2018-11-15:
```
Benchmark_Add/big-6  	100000000	        21.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/uint256-6         	300000000	         5.27 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/big-6             	100000000	        22.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/uint256-6         	300000000	         4.38 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/uint256_of-6      	300000000	         4.75 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/big-6             	10000000	       149 ns/op	     128 B/op	       2 allocs/op
Benchmark_Mul/uint256-6         	30000000	        56.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Square/big-6          	10000000	       144 ns/op	     128 B/op	       2 allocs/op
Benchmark_Square/uint256-6      	30000000	        46.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/big-6             	100000000	        13.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_And/uint256-6         	2000000000	         1.89 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/big-6              	100000000	        17.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/uint256-6          	2000000000	         1.88 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/big-6             	100000000	        18.1 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/uint256-6         	2000000000	         1.88 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/big-6             	200000000	         7.86 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/uint256-6         	500000000	         3.96 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/big/n_eq_0-6      	20000000	       100 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_192-6    	20000000	        93.1 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_128-6    	20000000	        93.1 ns/op	     128 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_64-6     	20000000	        89.9 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/big/n_gt_0-6      	20000000	        87.8 ns/op	     112 B/op	       2 allocs/op
Benchmark_Lsh/uint256/n_eq_0-6  	500000000	         3.71 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_192-6         	300000000	         4.23 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_128-6         	300000000	         5.64 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_64-6          	200000000	         8.08 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_0-6           	200000000	         9.77 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/big/n_eq_0-6               	20000000	        88.7 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_192-6             	20000000	        78.7 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_128-6             	20000000	        79.0 ns/op	      80 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_64-6              	20000000	        81.3 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/big/n_gt_0-6               	20000000	        82.6 ns/op	      96 B/op	       2 allocs/op
Benchmark_Rsh/uint256/n_eq_0-6           	500000000	         3.60 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_192-6         	300000000	         4.34 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_128-6         	300000000	         5.62 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_64-6          	200000000	         7.84 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_0-6           	200000000	         9.63 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/big-6                	   50000	     26358 ns/op	   18224 B/op	     191 allocs/op
Benchmark_Exp/large/uint256-6            	  100000	     17744 ns/op	      32 B/op	       1 allocs/op
Benchmark_Exp/small/big-6                	  200000	      7530 ns/op	    7472 B/op	      79 allocs/op
Benchmark_Exp/small/uint256-6            	 1000000	      1541 ns/op	      32 B/op	       1 allocs/op
Benchmark_Div/large/big-6                	 5000000	       311 ns/op	     176 B/op	       3 allocs/op
Benchmark_Div/large/uint256-6            	 5000000	       241 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/small/big-6                	10000000	       142 ns/op	     128 B/op	       3 allocs/op
Benchmark_Div/small/uint256-6            	100000000	        14.5 ns/op	       0 B/op	       0 allocs/op
Benchmark_MulMod/large/big-6             	 3000000	       536 ns/op	     320 B/op	       4 allocs/op
Benchmark_MulMod/large/uint256-6         	 1000000	      1245 ns/op	     608 B/op	      11 allocs/op
Benchmark_MulMod/small/big-6             	10000000	       165 ns/op	     128 B/op	       3 allocs/op
Benchmark_MulMod/small/uint256-6         	30000000	        54.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mod/large/big-6                	10000000	       133 ns/op	      48 B/op	       1 allocs/op
Benchmark_Mod/large/uint256-6            	10000000	       240 ns/op	     100 B/op	       3 allocs/op
Benchmark_Mod/small/big-6                	20000000	        63.4 ns/op	      48 B/op	       1 allocs/op
Benchmark_Mod/small/uint256-6            	100000000	        16.3 ns/op	       0 B/op	       0 allocs/op
Benchmark_SDiv/large/big-6               	 3000000	       517 ns/op	     352 B/op	       6 allocs/op
Benchmark_SDiv/large/uint256-6           	 5000000	       255 ns/op	     128 B/op	       3 allocs/op
Benchmark_AddMod/large/big-6             	 5000000	       263 ns/op	     128 B/op	       2 allocs/op
Benchmark_AddMod/large/uint256-6         	30000000	        56.5 ns/op	       0 B/op	       0 allocs/op
Benchmark_AddMod/small/big-6             	20000000	        95.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_AddMod/small/uint256-6         	50000000	        24.9 ns/op	       0 B/op	       0 allocs/op

```

The fixed lib wins over big in most cases, with a few exceptions: 

- `MulMod` on large numbers is slower on `uint256` by ~2x, since they wrap `big.Int` when the multiplication 
would overflow `256` bits.
- `Mod` on large numbers is slower on `uint256` by ~2x. 

Also, the tests are a bit lacking, specifically tests to ensure that function on the form `func (z *Int) Foo(a,b *Int)` 
* does not 
modify `a` or `b`, if `z.Foo(a,b)` and `z != (a,b)`
* works correctly in both cases `a.Foo(a,b)` and `b.Foo(a,b)`

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
