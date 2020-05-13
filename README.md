# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereum/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `uint256` being this library. 

### Current status

- As of 2020-03-18, `uint256` wins over big in every single case, often with orders of magnitude.
- And as of release `0.1.0`, the `uint256` library is alloc-free. 
- With the `1.0.0` release, it also has `100%` test coverage. 
 
### Conversion from/to `big.Int`

```
BenchmarkSetFromBig/1word-6        	252819142	         5.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/2words-6       	246641738	         5.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/3words-6       	238986948	         5.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/4words-6       	235735474	         5.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkSetFromBig/overflow-6     	233997644	         5.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkToBig/1word-6             	16905244	        76.4 ns/op	      64 B/op	       2 allocs/op
BenchmarkToBig/2words-6            	16532232	        75.3 ns/op	      64 B/op	       2 allocs/op
BenchmarkToBig/3words-6            	16193755	        73.7 ns/op	      64 B/op	       2 allocs/op
BenchmarkToBig/4words-6            	16226127	        76.5 ns/op	      64 B/op	       2 allocs/op
```
### Math operations

`uint256`:
```
Benchmark_Add/single/uint256-6     	629591630	         1.93 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/single/uint256-6     	546594793	         2.19 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/single/uint256_of-6  	569249097	         2.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkMul/single/uint256-6      	121926037	         8.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkSquare/single/uint256-6   	175970640	         6.86 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/uint256-6      	  325039	      3698 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/small/uint256-6      	 3832027	       328 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/small/uint256-6       	86487969	        14.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod64/uint256-6       	18414321	        66.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod128/uint256-6      	 9474686	       120 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod192/uint256-6      	11031271	       107 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv/mod256/uint256-6      	13804608	        87.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/small/uint256-6       	83189749	        14.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod64/uint256-6       	17734662	        67.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod128/uint256-6      	10050982	       119 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod192/uint256-6      	10477886	       112 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/mod256/uint256-6      	13230618	        92.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/small/uint256-6    	64955443	        19.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod64/uint256-6    	15589430	        79.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod128/uint256-6   	 8173804	       138 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod192/uint256-6   	 9115947	       133 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddMod/mod256/uint256-6   	10549335	       119 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/small/uint256-6    	29009958	        41.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod64/uint256-6    	10366879	       115 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod128/uint256-6   	 5339348	       256 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod192/uint256-6   	 5322434	       231 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-6   	 5229891	       229 ns/op	       0 B/op	       0 allocs/op
Benchmark_SDiv/large/uint256-6     	10094680	       122 ns/op	       0 B/op	       0 allocs/op
```
vs `big.Int`
```
Benchmark_Add/single/big-6         	48557786	        21.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/single/big-6         	55511476	        21.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMul/single/big-6          	17177685	        72.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSquare/single/big-6       	18434874	        67.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_Exp/large/big-6          	   30236	     34919 ns/op	   18144 B/op	     189 allocs/op
Benchmark_Exp/small/big-6          	  138970	      8290 ns/op	    7392 B/op	      77 allocs/op
BenchmarkDiv/small/big-6           	20982415	        55.9 ns/op	       8 B/op	       1 allocs/op
BenchmarkDiv/mod64/big-6           	 7928998	       168 ns/op	       8 B/op	       1 allocs/op
BenchmarkDiv/mod128/big-6          	 3858378	       312 ns/op	      80 B/op	       1 allocs/op
BenchmarkDiv/mod192/big-6          	 4733641	       251 ns/op	      80 B/op	       1 allocs/op
BenchmarkDiv/mod256/big-6          	 5928298	       204 ns/op	      80 B/op	       1 allocs/op
BenchmarkMod/small/big-6           	21673300	        63.5 ns/op	       8 B/op	       1 allocs/op
BenchmarkMod/mod64/big-6           	 7451844	       162 ns/op	      64 B/op	       1 allocs/op
BenchmarkMod/mod128/big-6          	 4138648	       308 ns/op	      64 B/op	       1 allocs/op
BenchmarkMod/mod192/big-6          	 4895036	       244 ns/op	      48 B/op	       1 allocs/op
BenchmarkMod/mod256/big-6          	 6530869	       208 ns/op	       8 B/op	       1 allocs/op
BenchmarkAddMod/small/big-6        	16872976	        73.4 ns/op	       8 B/op	       1 allocs/op
BenchmarkAddMod/mod64/big-6        	 6000974	       204 ns/op	      77 B/op	       1 allocs/op
BenchmarkAddMod/mod128/big-6       	 3480322	       355 ns/op	      64 B/op	       1 allocs/op
BenchmarkAddMod/mod192/big-6       	 3627595	       309 ns/op	      61 B/op	       1 allocs/op
BenchmarkAddMod/mod256/big-6       	 4768058	       247 ns/op	      40 B/op	       1 allocs/op
BenchmarkMulMod/small/big-6        	15861763	        77.0 ns/op	       8 B/op	       1 allocs/op
BenchmarkMulMod/mod64/big-6        	 3656503	       330 ns/op	      96 B/op	       1 allocs/op
BenchmarkMulMod/mod128/big-6       	 2104412	       577 ns/op	      96 B/op	       1 allocs/op
BenchmarkMulMod/mod192/big-6       	 2198616	       534 ns/op	      80 B/op	       1 allocs/op
BenchmarkMulMod/mod256/big-6       	 2492029	       536 ns/op	      80 B/op	       1 allocs/op
Benchmark_SDiv/large/big-6         	 2056039	       938 ns/op	     312 B/op	       6 allocs/op
```

### Boolean logic
`uint256`
```
Benchmark_And/single/uint256-6     	590464702	         1.93 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/single/uint256-6      	629714779	         2.06 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/single/uint256-6     	625039024	         1.94 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/single/uint256-6     	408978940	         2.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/large/uint256-6        	388216880	         3.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/small/uint256-6        	388475845	         3.37 ns/op	       0 B/op	       0 allocs/op

```
vs `big.Int`
```
Benchmark_And/single/big-6         	86528250	        14.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_Or/single/big-6          	69813108	        17.7 ns/op	       0 B/op	       0 allocs/op
Benchmark_Xor/single/big-6         	66239239	        18.2 ns/op	       0 B/op	       0 allocs/op
Benchmark_Cmp/single/big-6         	160121772	         7.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/large/big-6            	84170188	        12.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkLt/small/big-6            	84369852	        13.2 ns/op	       0 B/op	       0 allocs/op

```

### Bitwise shifts

`uint256`:
```
Benchmark_Lsh/n_eq_0/uint256-6     	260736044	         4.61 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_192/uint256-6   	208642209	         6.60 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_128/uint256-6   	185939301	         6.29 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_64/uint256-6    	134935269	         8.84 ns/op	       0 B/op	       0 allocs/op
Benchmark_Lsh/n_gt_0/uint256-6     	100000000	        10.6 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_eq_0/uint256-6     	267984342	         4.47 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_192/uint256-6   	226595919	         5.46 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_128/uint256-6   	174888495	         7.03 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_64/uint256-6    	135138018	         9.03 ns/op	       0 B/op	       0 allocs/op
Benchmark_Rsh/n_gt_0/uint256-6     	100000000	        10.8 ns/op	       0 B/op	       0 allocs/op
```
vs `big.Int`:
```
Benchmark_Lsh/n_eq_0/big-6         	22462352	        52.2 ns/op	      64 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_192/big-6       	18483247	        69.3 ns/op	      96 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_128/big-6       	17845354	        67.6 ns/op	      96 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_64/big-6        	19200085	        66.7 ns/op	      80 B/op	       1 allocs/op
Benchmark_Lsh/n_gt_0/big-6         	18988237	        63.1 ns/op	      80 B/op	       1 allocs/op
Benchmark_Rsh/n_eq_0/big-6         	22768628	        54.5 ns/op	      64 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_192/big-6       	29468035	        38.1 ns/op	       8 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_128/big-6       	23145200	        54.2 ns/op	      48 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_64/big-6        	19201986	        86.3 ns/op	      64 B/op	       1 allocs/op
Benchmark_Rsh/n_gt_0/big-6         	20327937	        62.3 ns/op	      64 B/op	       1 allocs/op
```
## Helping out

If you're interested in low-level algorithms and/or doing optimizations for shaving off nanoseconds, then this is certainly for you!

### Implementation work

Choose an operation, and optimize the s**t out of it!

A few rules, though, to help your PR get approved:

- Do not optimize for 'best-case'/'most common case' at the expense of worst-case. 
- We'll hold off on go assembly for a while, until the algos and interfaces are finished in a 'good enough' first version. After that, it's assembly time. 

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
