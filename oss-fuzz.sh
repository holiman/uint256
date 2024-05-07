#!/bin/bash -eu

function compile_fuzzer() {
  package=$1
  function=$2
  fuzzer=$3
  file=$4

  path=$GOPATH/src/$package

  echo "Building $fuzzer"
  cd $path

  # Install build dependencies
  go mod tidy
  go get github.com/holiman/gofuzz-shim/testing

	if [[ $SANITIZER == *coverage* ]]; then
		coverbuild $path $function $fuzzer $coverpkg
	else
	  gofuzz-shim --func $function --package $package -f $file -o $fuzzer.a
		$CXX $CXXFLAGS $LIB_FUZZING_ENGINE $fuzzer.a -o $OUT/$fuzzer
	fi

  ## Check if there exists a seed corpus file
  corpusfile="${path}/testdata/${fuzzer}_seed_corpus.zip"
  if [ -f $corpusfile ]
  then
    cp $corpusfile $OUT/
    echo "Found seed corpus: $corpusfile"
  fi
  cd -
}

go install github.com/holiman/gofuzz-shim@latest

repo=$GOPATH/src/github.com/holiman/uint256

compile_fuzzer github.com/holiman/uint256  FuzzUnaryOperations fuzzUnary $repo/unary_test.go
compile_fuzzer github.com/holiman/uint256  FuzzBinaryOperations fuzzBinary $repo/binary_test.go
compile_fuzzer github.com/holiman/uint256  FuzzCompareOperations fuzzCompare $repo/binary_test.go
compile_fuzzer github.com/holiman/uint256  FuzzTernaryOperations fuzzTernary $repo/ternary_test.go
compile_fuzzer github.com/holiman/uint256  FuzzSetString fuzzSetString $repo/uint256_test.go
