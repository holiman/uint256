package stack

import "github.com/holiman/uint256"

// 0x0 range - arithmetic ops.
const (
	STOP byte = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

// 0x10 range - comparison ops.
const (
	LT byte = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR

	SHA3 = 0x20
)
const (
	ADDRESS     = byte(0x30)
	ORIGIN      = byte(0x32)
	CALLER      = byte(0x33)
	CALLVALUE   = byte(0x34)
	GASPRICE    = byte(0x3a)
	COINBASE    = byte(0x41)
	TIMESTAMP   = byte(0x42)
	NUMBER      = byte(0x43)
	DIFFICULTY  = byte(0x44)
	GASLIMIT    = byte(0x45)
	CHAINID     = byte(0x46)
	SELFBALANCE = byte(0x47)
)

const (
	POP byte = 0x50 + iota
	MLOAD
	MSTORE
	MSTORE8
	SLOAD
	SSTORE
	JUMP
	JUMPI
	PC
	MSIZE
	GAS
	JUMPDEST
)

// 0x60 range.
// 0x60 range.
const (
	PUSH1 byte = 0x60 + iota
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
	DUP1
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
	SWAP1
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

type stackOp func(*StackMachine)

var stackFuncs [256]stackOp

func init() {

	stackFuncs[ADD] = (*StackMachine).opAdd
	stackFuncs[MUL] = (*StackMachine).opMul
	stackFuncs[SUB] = (*StackMachine).opSub
	stackFuncs[DIV] = (*StackMachine).opDiv
	stackFuncs[SDIV] = (*StackMachine).opSdiv
	stackFuncs[MOD] = (*StackMachine).opMod
	stackFuncs[SMOD] = (*StackMachine).opSMod
	stackFuncs[ADDMOD] = (*StackMachine).opAddmod
	stackFuncs[MULMOD] = (*StackMachine).opMulmod
	stackFuncs[EXP] = (*StackMachine).opExp
	stackFuncs[SIGNEXTEND] = (*StackMachine).opSignExtend

	stackFuncs[ORIGIN] = (*StackMachine).opOrigin
	stackFuncs[GASPRICE] = (*StackMachine).opGasprice
	stackFuncs[COINBASE] = (*StackMachine).opCoinbase
	stackFuncs[TIMESTAMP] = (*StackMachine).opTimestamp
	stackFuncs[NUMBER] = (*StackMachine).opNumber
	stackFuncs[DIFFICULTY] = (*StackMachine).opDifficulty
	stackFuncs[CHAINID] = (*StackMachine).opChainId
}

func (machine *StackMachine) DispatchSimple(op byte) (valid bool) {
	if fn := stackFuncs[op]; fn != nil {
		machine.fn()
		return true
	}
	if op >= PUSH1 && op <= PUSH32 {
		l := 1 + op - PUSH1
		if len(code) < int(l) {
			l = byte(len(code))
		}
		stack.PushBytes(code[:l])
		return true
	} else if op >= DUP1 && op <= DUP16 {
		stack.Dup(int(op - DUP1))
		return true
	} else if op >= SWAP1 && op <= SWAP16 {
		stack.Swap(int(op - SWAP1))
		return true
	}
	return false
}

// Call is the entry-point for execution, and is used for the initial outer call
// as well as any internal sub-calls, such as DELEGATECALL and STATICCALL
func (machine *StackMachine) Call(code []byte, address, caller [20]byte, value uint256.Int, readOnly bool, gas uint64) {
	if len(code) == 0 {
		return
	}
	// New context
	machine.NewContext()
	defer machine.DropContext()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This makes also sure that the readOnly flag isn't removed for child calls.
	if readOnly && !machine.static {
		machine.static = true
		defer func() { machine.static = false }()
	}
	var (
		pc           uint64
		opStaticCost uint64
	)

	for !abort {
		op = code[pc]

		// An opcode with a negative static gas value is not valid
		if cost := machine.staticGasCost[op]; cost < 0 {
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
		} else {
			opStaticCost = cost
		}
		// Validate stack
		if sLen := machine.StackDepth(); sLen < operation.minStack {
			return nil, fmt.Errorf("stack underflow (%d <=> %d)", sLen, operation.minStack)
		} else if sLen > operation.maxStack {
			return nil, fmt.Errorf("stack limit reached %d (%d)", sLen, operation.maxStack)
		}
		// If the operation is valid, enforce write restrictions
		if machine.static {
			// static mode can only be set if Byzantium is enabled, so no need
			// to check that explicitly
			// We have to make sure that no
			// state-modifying operation is performed. The 3rd stack item
			// for a call operation is the value. Transferring value from one
			// account to the others means the state is modified and should also
			// return with an error.
			//if operation.writes || (op == CALL && stack.Back(2).Sign() != 0) {
			//	return nil, errWriteProtection
			//}
		}
		if opStaticCost > gas {
			return nil, ErrOutOfGas
		}
		gas -= opStaticCost
		if machine.DispatchSimple(op) {
			continue
		}
		// All 'simple' opcodes are done. Remaining are calls, tx-context stackbuffer,
		// and a few otheres
		machine.DispatchComplex(op)

		// TODO: calc dynamic gas
		// TODO: memory expansion

	}

}

func (machine *StackMachine) memoryCall(inOff, inSize, retOff, retSize *uint256.Int) (uint64, bool) {
	x, overflow := calcMemSize64(inOff, inSize)
	if overflow {
		return 0, true
	}
	y, overflow := calcMemSize64(retOff, retSize)
	if overflow {
		return 0, true
	}
	if x > y {
		return x, false
	}
	return y, false
}

func (machine *StackMachine) dynamicGasCall(value *uint256.Int, destination [20]byte) {
	var (
		gasCost        uint64
		transfersValue = !value.IsZero()
	)
	// Determine how much gas is needed for the memory expansion (if any)
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, 0, err
	}
	// Determine extra gas that needs to be paid for creating a new account
	if isEip158 {
		if transfersValue && machine.StateDB.Empty(destination) {
			gasCost += params.CallNewAccountGas
		}
	} else if !machine.StateDB.Exist(destination) {
		gasCost += params.CallNewAccountGas
	}
	// Extra fee for transferring ether
	if transfersValue {
		gasCost += params.CallValueTransferGas
	}
	// Add the two factors
	var overflow bool
	if gasCost, overflow = math.SafeAdd(gasCost, memoryGas); overflow {
		return 0, 0, errGasUintOverflow
	}
	// Apply the 63/64:ths rule, to determine how much gas is actually
	// passed along to the callee
	sentGas, err = callGas(isEip150, availableGas, gasCost, desiredGas)
	if err != nil {
		return 0, 0, err
	}
	// Check if we have sufficient gas
	if gasCost, overflow = math.SafeAdd(availableGas, sentGas); overflow {
		return 0, 0, errGasUintOverflow
	}
	return gasCost, sentGas, nil
}

func (machine *StackMachine) opCall(availableGas uint64) error {

	var (
		desiredGas  = uint256.NewInt()
		destination [20]byte
		value       = uint256.NewInt()
		inOffset    = uint256.NewInt()
		inSize      = uint256.NewInt()
		retOffset   = uint256.NewInt()
		retSize     = uint256.NewInt()
	)
	machine.PopUint(desiredGas)
	machine.PopBytes20(destination)
	machine.PopUint(value)
	machine.PopUints(inOffset, inSize)
	machine.PopUints(retOffset, retSize)

	// We can do the memory expansion calculation and dynamic gas calc here
	memSize, overflow := machine.memoryCall(inOffset, inSize, retOffset, retSize)
	if overflow {
		return 0, errGasUintOverflow
	}
	dynamicGasCost, calleeGas, err := machine.dynamicGasCall(value, destination, memSize)
	if err != nil {
		return 0, err
	}
	if availableGas < dynamicGasCost {
		return 0, ErrOutOfGas
	}
	availableGas -= dynamicGasCost
	// At this point, we know there is sufficient gas, so it's time to
	// expand the memory
	if memSize > 0 {
		machine.callCtx.memory.Resize(memorySize)
	}
	// Execute the CALL
	ret, returnedGas, err := machine.Call(code, address, caller, value, false, calleeGas)
	if err != nil {
		machine.PushZero()
	} else {
		machine.PushBytes(oneByte32)
	}
	if err == nil || err == errExecutionReverted {
		machine.memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
		machine.returnData = ret
	}
	return availableGas + returnedGas, nil
}
