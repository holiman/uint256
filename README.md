# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereum/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `uint256` being this library. 

As of 2020-03-18, the fixed lib wins over big in every single case, often with orders of magnitude.
As of release `0.1.0`, the `uint256` library is alloc-free!
 
### Conversion from/to `big.Int`

```
BenchmarkSetFromBig/1word-6        	259504869	         4.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/2words-6       	243296958	         5.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/3words-6       	227551600	         5.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/4words-6       	246922267	         5.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/overflow-6     	238420483	         5.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkToBig/1word-6             	14952933	        71.1 ns/op	      64 B/op	       2 allocs/op
BenchmarkToBig/2words-6            	17000169	        72.5 ns/op	      64 B/op	       2 allocs/op
BenchmarkToBig/3words-6            	14798148	        72.2 ns/op	      64 B/op	       2 allocs/op
BenchmarkToBig/4words-6            	15954582	        69.9 ns/op	      64 B/op	       2 allocs/op

```
### Math operations

`uint256`:
```
enchmark_Add/single/uint256-6     	633472188	         1.89 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/single/uint256-6     	568273075	         2.22 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/single/uint256_of-6  	567779286	         2.11 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/single/uint256-6     	62855088	        17.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_Square/single/uint256-6  	64012897	        18.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/uint256-6      	  153199	      7749 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/small/uint256-6      	 1808989	       663 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod64/uint256-6       	18118854	        64.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod128/uint256-6      	10897378	       109 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod192/uint256-6      	10550451	        99.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod256/uint256-6      	13321220	        84.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/small/uint256-6       	78706082	        14.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod64/uint256-6       	17726313	        67.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod128/uint256-6      	10611421	       112 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod192/uint256-6      	11801432	       101 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod256/uint256-6      	13408267	        86.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/small/uint256-6    	63672259	        18.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod64/uint256-6    	13364508	        76.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod128/uint256-6   	 9374313	       128 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod192/uint256-6   	 9662953	       121 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod256/uint256-6   	11169826	       108 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/small/uint256-6    	26907108	        44.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod64/uint256-6    	10069104	       120 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod128/uint256-6   	 5551214	       215 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod192/uint256-6   	 5677230	       209 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-6   	 5793927	       204 ns/op	       0 B/op	       0 allocs/op
Benchmark_SDiv/large/uint256-6     	11370470	       105 ns/op	       0 B/op	       0 allocs/op

```
vs `big.Int`
```
Benchmark_Add/single/big-6         	55087052	        21.5 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/single/big-6         	55020072	        21.4 ns/op	       0 B/op	       0 allocs/op
Benchmark_Mul/single/big-6         	 9587319	       119 ns/op	      96 B/op	       1 allocs/op
Benchmark_Square/single/big-6      	 9989365	       128 ns/op	      96 B/op	       1 allocs/op
Benchmark_Exp/large/big-6          	   44059	     27808 ns/op	   18264 B/op	     192 allocs/op
Benchmark_Exp/small/big-6          	  132482	      8094 ns/op	    7512 B/op	      80 allocs/op
BenchmarkDiv/small/big-6           	20346138	        56.6 ns/op	       8 B/op	       1 allocs/op
BenchmarkDiv/mod64/big-6           	 8671370	       139 ns/op	       8 B/op	       1 allocs/op
BenchmarkDiv/mod128/big-6          	 3985227	       301 ns/op	      80 B/op	       1 allocs/op
BenchmarkDiv/mod192/big-6          	 4860343	       249 ns/op	      80 B/op	       1 allocs/op
BenchmarkDiv/mod256/big-6          	 4927914	       256 ns/op	      80 B/op	       1 allocs/op
BenchmarkMod/small/big-6           	23049183	        51.9 ns/op	       8 B/op	       1 allocs/op
BenchmarkMod/mod64/big-6           	 7158260	       162 ns/op	      64 B/op	       1 allocs/op
BenchmarkMod/mod128/big-6          	 3976514	       290 ns/op	      64 B/op	       1 allocs/op
BenchmarkMod/mod192/big-6          	 4926200	       312 ns/op	      48 B/op	       1 allocs/op
BenchmarkMod/mod256/big-6          	 6534632	       180 ns/op	       8 B/op	       1 allocs/op
BenchmarkAddMod/small/big-6        	16317876	        71.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkAddMod/mod64/big-6        	 5799704	       201 ns/op	      77 B/op	       1 allocs/op
BenchmarkAddMod/mod128/big-6       	 3470521	       415 ns/op	      64 B/op	       1 allocs/op
BenchmarkAddMod/mod192/big-6       	 4080316	       295 ns/op	      61 B/op	       1 allocs/op
BenchmarkAddMod/mod256/big-6       	 4928460	       239 ns/op	      40 B/op	       1 allocs/op
BenchmarkMulMod/small/big-6        	15940292	        73.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkMulMod/mod64/big-6        	 3570828	       331 ns/op	      96 B/op	       1 allocs/op
BenchmarkMulMod/mod128/big-6       	 2047695	       597 ns/op	      96 B/op	       1 allocs/op
BenchmarkMulMod/mod192/big-6       	 2291578	       514 ns/op	      80 B/op	       1 allocs/op
BenchmarkMulMod/mod256/big-6       	 2463007	       509 ns/op	      80 B/op	       1 allocs/op
Benchmark_SDiv/large/big-6         	 2212314	       544 ns/op	     248 B/op	       5 allocs/op

```

### Boolean logic
`uint256`
```
Benchmark_And/single/uint256-6     	534546528	         2.16 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/single/uint256-6      	620022866	         1.90 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/single/uint256-6     	553041696	         2.15 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/single/uint256-6     	433571565	         2.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/large/uint256-6        	407432343	         2.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/small/uint256-6        	361346886	         2.95 ns/op	       0 B/op	       0 allocs/op

```
vs `big.Int`
```
Benchmark_And/single/big-6         	74969566	        15.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/single/big-6          	69831138	        18.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/single/big-6         	65390836	        17.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/single/big-6         	165527830	         7.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/large/big-6            	85042539	        13.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/small/big-6            	93413899	        12.8 ns/op	       0 B/op	       0 allocs/op

```

### Bitwise shifts

`uint256`:
```
Benchmark_Lsh/n_eq_0/uint256-6     	275789586	         4.36 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_192/uint256-6   	228063312	         5.30 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_128/uint256-6   	204222584	         5.81 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_64/uint256-6    	144790902	         8.27 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_0/uint256-6     	122010325	         9.85 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_eq_0/uint256-6     	272249740	         4.36 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_192/uint256-6   	240133351	         4.96 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_128/uint256-6   	180698709	         6.65 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_64/uint256-6    	142424036	         8.70 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_0/uint256-6     	120315340	         9.94 ns/op	       0 B/op	       0 allocs/op
```
vs `big.Int`:
```
Benchmark_Lsh/n_eq_0/big-6         	19975887	        51.0 ns/op	      64 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_192/big-6       	17185830	        71.5 ns/op	      96 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_128/big-6       	16629559	        65.7 ns/op	      96 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_64/big-6        	17201168	        62.9 ns/op	      80 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_0/big-6         	16718604	        60.9 ns/op	      80 B/op	       1 allocs/op
Benchmark_Rsh/n_eq_0/big-6         	19647315	        52.7 ns/op	      64 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_192/big-6       	29589792	        37.4 ns/op	       8 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_128/big-6       	18169792	        74.5 ns/op	      48 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_64/big-6        	21416439	        55.2 ns/op	      64 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_0/big-6         	17568159	        58.7 ns/op	      64 B/op	       1 allocs/op
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
