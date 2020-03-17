# Fixed size math

This is a library specialized at replacing the big.Int library for math based on 256-bit types. This is meant for use in [go-ethereum](https://github.com/ethereum/go-ethereum) eventually, once it's deemed fast, stable and secure enough. 

## Benchmarks

Current benchmarks, with tests ending with `big` being the standard `big.Int` library, and `uint256` being this library. 

As of 2020-03-17, the fixed lib wins over big in every single case, often with orders of magnitude.
 
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
```
Benchmark_Add/big-6                     48508614                24.8 ns/op             0 B/op          0 allocs/op
Benchmark_Add/uint256-6                 301348272                3.95 ns/op            0 B/op          0 allocs/op
Benchmark_Sub/big-6                     52758924                26.5 ns/op             0 B/op          0 allocs/op
Benchmark_Sub/uint256-6                 312972354                3.95 ns/op            0 B/op          0 allocs/op
Benchmark_Sub/uint256_of-6              541212834                2.30 ns/op            0 B/op          0 allocs/op
Benchmark_Mul/big-6                      9169669               131 ns/op              96 B/op          1 allocs/op
Benchmark_Mul/uint256-6                 55652130                22.6 ns/op             0 B/op          0 allocs/op
Benchmark_Square/big-6                   9017666               128 ns/op              96 B/op          1 allocs/op
Benchmark_Square/uint256-6              54524514                22.8 ns/op             0 B/op          0 allocs/op
Benchmark_Exp/large/big-6                  42733             30038 ns/op           18264 B/op        192 allocs/op
Benchmark_Exp/large/uint256-6             120678              9212 ns/op               0 B/op          0 allocs/op
Benchmark_Exp/small/big-6                 137761              8947 ns/op            7512 B/op         80 allocs/op
Benchmark_Exp/small/uint256-6            1496252               790 ns/op               0 B/op          0 allocs/op
Benchmark_Div/large/big-6                3631567               338 ns/op             144 B/op          2 allocs/op
Benchmark_Div/large/uint256-6            6181081               199 ns/op              64 B/op          2 allocs/op
Benchmark_Div/small/big-6               12047379                98.3 ns/op            16 B/op          2 allocs/op
Benchmark_Div/small/uint256-6           65821525                18.8 ns/op             0 B/op          0 allocs/op
Benchmark_MulMod/large/big-6             2221353               724 ns/op             176 B/op          2 allocs/op
Benchmark_MulMod/large/uint256-6         3551917               338 ns/op              80 B/op          2 allocs/op
Benchmark_MulMod/small/big-6             8594644               132 ns/op              56 B/op          2 allocs/op
Benchmark_MulMod/small/uint256-6        20614293                56.7 ns/op             0 B/op          0 allocs/op
Benchmark_Mod/large/big-6                6542557               181 ns/op               8 B/op          1 allocs/op
Benchmark_Mod/large/uint256-6           10740315               113 ns/op              40 B/op          2 allocs/op
Benchmark_Mod/small/big-6               20946499                58.0 ns/op             8 B/op          1 allocs/op
Benchmark_Mod/small/uint256-6           66800676                19.2 ns/op             0 B/op          0 allocs/op
Benchmark_SDiv/large/big-6               1985832               636 ns/op             248 B/op          5 allocs/op
Benchmark_SDiv/large/uint256-6           5712229               210 ns/op              64 B/op          2 allocs/op
Benchmark_AddMod/large/big-6             4206482               287 ns/op              48 B/op          1 allocs/op
Benchmark_AddMod/large/uint256-6        22825184                52.9 ns/op             0 B/op          0 allocs/op
Benchmark_AddMod/small/big-6            13530213                88.4 ns/op             8 B/op          1 allocs/op
Benchmark_AddMod/small/uint256-6        51714418                23.1 ns/op             0 B/op          0 allocs/op
```

### Boolean logic
```
Benchmark_And/big-6                     59583842                18.2 ns/op             0 B/op          0 allocs/op
Benchmark_And/uint256-6                 518743656                2.58 ns/op            0 B/op          0 allocs/op
Benchmark_Or/big-6                      62413881                19.4 ns/op             0 B/op          0 allocs/op
Benchmark_Or/uint256-6                  503491710                2.42 ns/op            0 B/op          0 allocs/op
Benchmark_Xor/big-6                     57397172                20.5 ns/op             0 B/op          0 allocs/op
Benchmark_Xor/uint256-6                 528165157                2.41 ns/op            0 B/op          0 allocs/op
Benchmark_Cmp/big-6                     143009628                8.32 ns/op            0 B/op          0 allocs/op
Benchmark_Cmp/uint256-6                 279421581                4.34 ns/op            0 B/op          0 allocs/op
```
### Bitwise shifts

```
Benchmark_Lsh/big/n_eq_0-6              20682451                75.0 ns/op            64 B/op          1 allocs/op
Benchmark_Lsh/big/n_gt_192-6            15043952                73.5 ns/op            96 B/op          1 allocs/op
Benchmark_Lsh/big/n_gt_128-6            16426294                71.8 ns/op            96 B/op          1 allocs/op
Benchmark_Lsh/big/n_gt_64-6             17304938                68.3 ns/op            80 B/op          1 allocs/op
Benchmark_Lsh/big/n_gt_0-6              16005706                72.4 ns/op            80 B/op          1 allocs/op
Benchmark_Lsh/uint256/n_eq_0-6          244153966                4.90 ns/op            0 B/op          0 allocs/op
Benchmark_Lsh/uint256/n_gt_192-6        208895947                5.94 ns/op            0 B/op          0 allocs/op
Benchmark_Lsh/uint256/n_gt_128-6        177531820                7.53 ns/op            0 B/op          0 allocs/op
Benchmark_Lsh/uint256/n_gt_64-6         127253976                9.41 ns/op            0 B/op          0 allocs/op
Benchmark_Lsh/uint256/n_gt_0-6          100000000               11.1 ns/op             0 B/op          0 allocs/op
Benchmark_Rsh/big/n_eq_0-6              20687982                58.1 ns/op            64 B/op          1 allocs/op
Benchmark_Rsh/big/n_gt_192-6            29531353                40.4 ns/op             8 B/op          1 allocs/op
Benchmark_Rsh/big/n_gt_128-6            21137516                57.1 ns/op            48 B/op          1 allocs/op
Benchmark_Rsh/big/n_gt_64-6             18891128                61.9 ns/op            64 B/op          1 allocs/op
Benchmark_Rsh/big/n_gt_0-6              17295828                63.0 ns/op            64 B/op          1 allocs/op
Benchmark_Rsh/uint256/n_eq_0-6          247489960                5.30 ns/op            0 B/op          0 allocs/op
Benchmark_Rsh/uint256/n_gt_192-6        211331108                5.52 ns/op            0 B/op          0 allocs/op
Benchmark_Rsh/uint256/n_gt_128-6        177311349                7.02 ns/op            0 B/op          0 allocs/op
Benchmark_Rsh/uint256/n_gt_64-6         125465119                9.36 ns/op            0 B/op          0 allocs/op
Benchmark_Rsh/uint256/n_gt_0-6          100000000               11.3 ns/op             0 B/op          0 allocs/op
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
