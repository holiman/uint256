package stack

import (
	"encoding/binary"
	"fmt"
	"github.com/holiman/uint256"
)

var (
	zeroByte12 = make([]byte, 12)
	zeroByte32 = make([]byte, 32)
	byteFalse  = zeroByte32
	byteTrue   = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	oneByte32 = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// CallContext represents the call context for an execution, meaning the
// address, memory and stack resulting from either a toplevel transaction or
// a CALL-type.
// It does not contain ephemeral things like pc
type CallContext struct {
	memory      *Memory
	stackHead   int // index of next free byte
	stackBottom int // index of first stack item

}

type StackMachine struct {

	// stackbuffer is the underlying global evm stackbuffer-stack, which becomes
	// partitioned into regions by the contexts. It is heavily reused across
	// contexts, transactions and even blocks.
	stackbuffer []byte
	// contexts is a stack of previous call-contexts
	contexts []*CallContext
	// callCtx is the current executing call-context
	callCtx *CallContext

	// allocated tracks the size of the underlying (global) stackbuffer
	allocated int
	// The ROM is where we keep per-block and per-tx constants, like ORIGIN and TIMESTAMP
	rom readonlyMem
	// whether we're in read-only (static) mode
	static bool
	// Some temp variables. These are used for arithmetic operations, and
	// are never leaked outside of the StackMachine
	x *uint256.Int
	y *uint256.Int
	z *uint256.Int
}

func New() *StackMachine {
	size := 1024
	return &StackMachine{
		// 1024B = room for 32 stack elements by default
		stackbuffer: make([]byte, size),
		allocated:   size,
		rom:         newReadonlyMem(),
		x:           uint256.NewInt(),
		y:           uint256.NewInt(),
		z:           uint256.NewInt(),
	}
}

// NewContext starts a new stack context
// A context is the CALL-context, meaning
// 1) A local stack
// 2) A local memory area
func (machine *StackMachine) NewContext() {
	ctx := &CallContext{
		memory:      NewMemory(),
		stackHead:   0,
		stackBottom: 0,
	}
	if cur := machine.callCtx; cur != nil {
		ctx.stackBottom = cur.stackHead
		ctx.stackHead = cur.stackHead
		machine.contexts = append(machine.contexts, cur)
	}
}

// DropContext drops the current context, restoring the previous context
func (machine *StackMachine) DropContext() {
	idx := len(machine.contexts) - 1
	machine.callCtx = machine.contexts[ixd]
	machine.contexts = machine.contexts[:idx]
}

// StackDepth returns the number of items in the current stack context
func (machine *StackMachine) StackDepth() {
	bottom := machine.callCtx.stackBottom
	head := machine.callCtx.stackHead
	return (head - bottom) / 32
}

func (machine *StackMachine) InitTx(origin, gasPrice *big.Int) {
	// TODO
	// - clear any contexts
	//   -- maybe add a panic if they were not properly
	// 		cleared by deferred dropContexts
	// - Maybe un-grow the stackbuffer area
	machine.rom.setTxConstants(origin, gasPrice)
}
func (machine *StackMachine) InitBlock(coinbase, timestamp, number,
	difficulty, chainid *big.Int) {
	// TODO
	// - clear any contexts
	//   -- maybe add a panic if they were not properly
	// 		cleared by deferred dropContexts
	// - Maybe un-grow the stackbuffer area
	machine.rom.setBlockConstants(coinbase, timestamp, number, difficulty, chainid)
}

// maybeGrow checks if the stack is large enough, otherwise it will grow a bit
func (machine *StackMachine) maybeGrow() {
	if machine.allocated < machine.head+32 {
		factor := len(machine.stackbuffer) / 2
		machine.allocated += factor
		newSlice := make([]byte, factor)
		machine.stackbuffer = append(machine.stackbuffer, newSlice...)
	}
}

// PushBytes copies the stackbuffer into the stack. It's safe to modify the argument
// after the function has returned.
// If the stackbuffer is smaller than 32 bytes, the stack will be zero-filled from the
// left.
// If the stackbuffer is larger than 32 bytes, it will be cropped from the left, leaving
// only the least significant 32 bytes.
func (machine *StackMachine) PushBytes(data []byte) {
	machine.maybeGrow()
	len := len(data)
	head := machine.callCtx.stackHead
	if len < 32 {
		// Clear the area
		copy(machine.stackbuffer[head:], zeroByte32)
		copy(machine.stackbuffer[head+(32-len):], data)
	} else {
		copy(machine.stackbuffer[head:], data[len-32:len])
	}
	machine.callCtx.stackHead += 32
}

// Dup duplicates the n:th item, placing the duplicate on the top stack position
// n == 0 means first item
func (machine *StackMachine) Dup(n int) {
	machine.maybeGrow()
	from := machine.callCtx.stackHead - 32 - 32*n
	copy(machine.stackbuffer[machine.callCtx.stackHead:], stack.stackbuffer[from:from+32])
	machine.callCtx.head += 32
}

// Swap swaps the n:th item with the head item in the stack
func (machine *StackMachine) Swap(n int) {
	// Swap doesn't increase the stack, but we use the 'future' space
	// as a swap area, so need to ensure we don't go out of bounds of the
	// allocated area
	machine.maybeGrow()
	head := machine.callCtx.stackHead
	src := head - 32*(n+1)
	// Copy n:th item to scratch-space
	copy(machine.stackbuffer[head:], machine.stackbuffer[src:src+32])
	// Copy head item to n:th place
	copy(machine.stackbuffer[src:], machine.stackbuffer[head-32:head])
	// Copy temp item to head place
	copy(machine.stackbuffer[head-32:head], machine.stackbuffer[head:head+32])
}

// PushBytes32 pushes 32 bytes of stackbuffer onto the stack
func (machine *StackMachine) PushBytes32(data [32]byte) {
	machine.maybeGrow()
	copy(machine.stackbuffer[stack.callContext.stackHead:], data[:])
	machine.callCtx.stackHead += 32
}

// PushBytes20 pushes 20 bytes of stackbuffer onto the stack
func (machine *StackMachine) PushBytes20(data [20]byte) {
	machine.maybeGrow()
	head := machine.callCtx.stackHead
	copy(machine.stackbuffer[head:], zeroByte12)
	copy(machine.stackbuffer[head+12:], data[:])
	machine.callCtx.stackHead += 32
}

// PushUint64 pushes a uint64 onto the stack
func (machine *StackMachine) PushUint64(x uint64) {
	machine.maybeGrow()
	head := machine.callCtx.stackHead
	copy(machine.stackbuffer[head:], zeroByte32[:24])
	binary.BigEndian.PutUint64(machine.stackbuffer[head+24:], x)
	machine.callCtx.stackHeadhead += 32
}

// PushUint pushes a uint256.Int onto the stack. It's safe to modify the
// argument after this method returns
func (machine *StackMachine) PushUint(u256 *uint256.Int) {
	machine.maybeGrow()
	// WriteToSlice writes the full 32 bytes, no need to clear space
	u256.WriteToSlice(machine.stackbuffer[machine.callCtx.stackHead:])
	machine.callCtx.stackHead += 32
}

// PushBool pushes a bool onto the stack
func (machine *StackMachine) PushBool(v bool) {
	machine.maybeGrow()
	if v {
		copy(machine.stackbuffer[machine.callCtx.stackHead:], byteTrue)
	} else {
		copy(machine.stackbuffer[machine.callCtx.stackHead:], byteFalse)
	}
	machine.callCtx.stackHead += 32
}

// PushZero pushes 32-bytes of zeroes onto the stack
func (machine *StackMachine) PushZero() {
	machine.maybeGrow()
	copy(machine.stackbuffer[machine.callCtx.stackHead:], zeroByte32)
}

// Pop removes the head element on the stack
func (machine *StackMachine) Pop() {
	machine.callCtx.stackHead -= 32
}

// PopUint pops the head element as a uint256.Int
func (machine *StackMachine) PopUint(u256 *uint256.Int) {
	h := machine.callCtx.stackHead
	u256.SetBytes(machine.stackbuffer[h-32 : h])
	machine.callCtx.stackHead -= 32
}

// PopUint pops two head elements as uint256.Ints
func (machine *StackMachine) PopUints(x, y *uint256.Int) {
	h := machine.callCtx.stackHead
	x.SetBytes(machine.stackbuffer[h-32 : h])
	y.SetBytes(machine.stackbuffer[h-64 : h-32])
	machine.callCtx.stackHead -= 64
}

// PopBytes32 pops 32 byte of stackbuffer off the stack and copies into stackbuffer
func (machine *StackMachine) PopBytes32(data [32]byte) [32]byte {
	h := machine.callContext.stackHead
	copy(data[:], machine.stackbuffer[h-32:h])
	machine.callContext.stackHead -= 32
	return data
}

func (machine *StackMachine) PrettyPrint() {
	fmt.Printf("StackMachine\n")
	fmt.Printf("  allocated: %d\n", machine.allocated)
	fmt.Printf("  contexts: %d\n", machine.contexts)
	for i := 0; i < machine.callCtx.stackHead; i += 32 {
		s := machine.stackbuffer[i : i+32]
		fmt.Printf("%02d  %x\n", i/32, s)
	}
}
