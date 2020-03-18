package stack

import (
	"github.com/holiman/uint256"
	"testing"
)

func TestSimpleLoop(t *testing.T) {
	// 0xfffff = 1048575 loops
	code := []byte{
		byte(PUSH3), 0x0f, 0xff, 0xff,
		byte(JUMPDEST), //  [ count ]
		byte(PUSH1), 1, // [count, 1]
		byte(SWAP1),    // [1, count]
		byte(SUB),      // [ count -1 ]
		byte(DUP1),     //  [ count -1 , count-1]
		byte(PUSH1), 4, // [count-1, count -1, label]
		byte(JUMPI), // [ 0 ]
		byte(STOP),
	}
	//tracer := vm.NewJSONLogger(nil, os.Stdout)
	//Execute(code, nil, &Config{
	//	EVMConfig: vm.Config{
	//		Debug:  true,
	//		Tracer: tracer,
	//	}})
	var (
		vm      = New()
		zero    = uint256.NewInt()
		address [20]byte
		caller  [20]byte
	)
	vm.InitBlock(zero, zero, zero, zero, zero)
	vm.InitTx(zero, zero)
	//for i := 0; i < b.N; i++ {
	vm.Call(code, address, caller, zero, false, 0xffffffff)
	//execute(code)
	//}
}

// BenchmarkSimpleLoop test a pretty simple loop which loops
// 1M (1 048 575) times.
// Takes about 200 ms in go-ethereum
// Takes about double that here :(
func BenchmarkSimpleLoop(b *testing.B) {
	// 0xfffff = 1048575 loops
	code := []byte{
		byte(PUSH3), 0x0f, 0xff, 0xff,
		byte(JUMPDEST), //  [ count ]
		byte(PUSH1), 1, // [count, 1]
		byte(SWAP1),    // [1, count]
		byte(SUB),      // [ count -1 ]
		byte(DUP1),     //  [ count -1 , count-1]
		byte(PUSH1), 4, // [count-1, count -1, label]
		byte(JUMPI), // [ 0 ]
		byte(STOP),
	}
	//tracer := vm.NewJSONLogger(nil, os.Stdout)
	//Execute(code, nil, &Config{
	//	EVMConfig: vm.Config{
	//		Debug:  true,
	//		Tracer: tracer,
	//	}})
	var (
		vm      = New()
		zero    = uint256.NewInt()
		address [20]byte
		caller  [20]byte
	)
	vm.InitBlock(zero, zero, zero, zero, zero)
	vm.InitTx(zero, zero)
	for i := 0; i < b.N; i++ {
		vm.Call(code, address, caller, zero, false, 0xffffffff)
	}
}
