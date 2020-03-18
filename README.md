# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereum/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `uint256` being this library. 

As of 2020-03-18, the fixed lib wins over big in every single case, often with orders of magnitude.
 
### Conversion from/to `big.Int`

```
BenchmarkSetFromBig/1word-6             258336367                4.76 ns/op            0 B/op          0 allocs/op
BenchmarkSetFromBig/2words-6            225051297                5.27 ns/op            0 B/op          0 allocs/op
BenchmarkSetFromBig/3words-6            200909203                5.34 ns/op            0 B/op          0 allocs/op
BenchmarkSetFromBig/4words-6            217516165                5.85 ns/op            0 B/op          0 allocs/op
BenchmarkSetFromBig/overflow-6          202799508                5.67 ns/op            0 B/op          0 allocs/op
BenchmarkToBig/1word-6                  16016079                77.2 ns/op            64 B/op          2 allocs/op
BenchmarkToBig/2words-6                 15877190                76.3 ns/op            64 B/op          2 allocs/op
BenchmarkToBig/3words-6                 15553672                77.8 ns/op            64 B/op          2 allocs/op
BenchmarkToBig/4words-6                 14967573                75.6 ns/op            64 B/op          2 allocs/op
```
### Math operations

`uint256`:
```
Benchmark_Add/uint256-6 	343733922	         3.48 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/uint256-6 	342691573	         3.49 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/uint256_of-6         	567302500	         2.12 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/uint256-6            	61316617	        19.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_Square/uint256-6         	59750362	        20.2 ns/op	       0 B/op	       0 allocs/op
```
vs `big.Int`
```
Benchmark_Add/big-6     	53215336	        23.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/big-6     	52027230	        22.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/big-6     	10366168	       116 ns/op	      96 B/op	       1 allocs/op
Benchmark_Square/big-6  	10606704	       116 ns/op	      96 B/op	       1 allocs/op
```

### Boolean logic
`uint256`
```
Benchmark_And/uint256-6            	571235724	         1.80 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/uint256-6             	630674871	         1.95 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/uint256-6            	629836861	         1.91 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/uint256-6            	267981819	         4.53 ns/op	       0 B/op	       0 allocs/op
```
vs `big.Int`
```
Benchmark_And/big-6     	74778183	        16.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/big-6      	69199390	        17.5 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/big-6     	63655377	        19.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/big-6     	148096443	         7.82 ns/op	       0 B/op	       0 allocs/op
```

### Bitwise shifts

`uint256`:
```
Benchmark_Lsh/uint256/n_eq_0-6     	276080395	         4.37 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_192-6   	236563666	         5.08 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_128-6   	193383686	         6.27 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_64-6    	134278740	         8.43 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/uint256/n_gt_0-6     	100000000	        10.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_eq_0-6     	273664575	         4.35 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_192-6   	237489868	         5.08 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_128-6   	184016202	         6.46 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_64-6    	139812172	         8.65 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/uint256/n_gt_0-6     	100000000	        10.9 ns/op	       0 B/op	       0 allocs/op
```
vs `big.Int`:
```
Benchmark_Lsh/big/n_eq_0-6         	22088068	        51.6 ns/op	      64 B/op	       1 allocs/op
Benchmark_Lsh/big/n_gt_192-6       	18656282	        64.3 ns/op	      96 B/op	       1 allocs/op
Benchmark_Lsh/big/n_gt_128-6       	18769039	        63.2 ns/op	      96 B/op	       1 allocs/op
Benchmark_Lsh/big/n_gt_64-6        	19760697	        68.4 ns/op	      80 B/op	       1 allocs/op
Benchmark_Lsh/big/n_gt_0-6         	20527872	        61.0 ns/op	      80 B/op	       1 allocs/op
Benchmark_Rsh/big/n_eq_0-6         	22072800	        52.7 ns/op	      64 B/op	       1 allocs/op
Benchmark_Rsh/big/n_gt_192-6       	32900992	        38.1 ns/op	       8 B/op	       1 allocs/op
Benchmark_Rsh/big/n_gt_128-6       	23719686	        73.0 ns/op	      48 B/op	       1 allocs/op
Benchmark_Rsh/big/n_gt_64-6        	22413007	        56.0 ns/op	      64 B/op	       1 allocs/op
Benchmark_Rsh/big/n_gt_0-6         	21509913	        57.0 ns/op	      64 B/op	       1 allocs/op
```
## Helping out

If you're interested in low-level algorithms and/or doing optimizations for shaving off nanoseconds, then this is certainly for you!

### Implementation work

Choose an operation, and optimize the s**t out of it!

A few rules, though, to help your PR get approved:

- Do not optimize for 'best-case'/'most common case' at the expense of worst-case. 
- We'll hold off on go assembly for a while, until the algos and interfaces are finished in a 'good enough' first version. After that, it's assembly time. 

### Testing
Also, any help in improving the test framework, e.g. by improving the random testing stuff is very highly appreciated. 

The tests are a bit lacking, specifically tests to ensure that function on the form `func (z *Int) Foo(a,b *Int)` 
* does not 
modify `a` or `b`, if `z.Foo(a,b)` and `z != (a,b)`
* works correctly in both cases `a.Foo(a,b)` and `b.Foo(a,b)`

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
