package stack

import (
	"errors"
	"fmt"
	"github.com/holiman/uint256"
)

type simpleOp func(*StackMachine)
type complexOp func(*StackMachine, *uint64, uint64) (uint64, error)

var simpleFuncs [256]simpleOp
var complexFuncs [256]complexOp
var staticGasCosts [256]uint64

var ErrOutOfGas = errors.New("Out of gas")

func init() {

	simpleFuncs[ADD] = (*StackMachine).opAdd
	simpleFuncs[MUL] = (*StackMachine).opMul
	simpleFuncs[SUB] = (*StackMachine).opSub
	simpleFuncs[DIV] = (*StackMachine).opDiv
	simpleFuncs[SDIV] = (*StackMachine).opSdiv
	simpleFuncs[MOD] = (*StackMachine).opMod
	simpleFuncs[SMOD] = (*StackMachine).opSMod
	simpleFuncs[ADDMOD] = (*StackMachine).opAddmod
	simpleFuncs[MULMOD] = (*StackMachine).opMulmod
	simpleFuncs[EXP] = (*StackMachine).opExp
	simpleFuncs[SIGNEXTEND] = (*StackMachine).opSignExtend

	simpleFuncs[ORIGIN] = (*StackMachine).opOrigin
	simpleFuncs[GASPRICE] = (*StackMachine).opGasprice
	simpleFuncs[COINBASE] = (*StackMachine).opCoinbase
	simpleFuncs[TIMESTAMP] = (*StackMachine).opTimestamp
	simpleFuncs[NUMBER] = (*StackMachine).opNumber
	simpleFuncs[DIFFICULTY] = (*StackMachine).opDifficulty
	simpleFuncs[CHAINID] = (*StackMachine).opChainId

	simpleFuncs[MSIZE] = (*StackMachine).opMsize
	simpleFuncs[JUMPDEST] = (*StackMachine).opJumpdest

	complexFuncs[MLOAD] = (*StackMachine).opTodo
	complexFuncs[MSTORE] = (*StackMachine).opTodo

	complexFuncs[STOP] = (*StackMachine).opTodo
	complexFuncs[ADDRESS] = (*StackMachine).opTodo
	complexFuncs[CALLER] = (*StackMachine).opTodo
	complexFuncs[CALLVALUE] = (*StackMachine).opTodo
	complexFuncs[CALLDATALOAD] = (*StackMachine).opTodo
	complexFuncs[CALLDATACOPY] = (*StackMachine).opTodo
	complexFuncs[CODESIZE] = (*StackMachine).opTodo
	complexFuncs[EXTCODESIZE] = (*StackMachine).opTodo
	complexFuncs[RETURNDATASIZE] = (*StackMachine).opTodo
	complexFuncs[EXTCODEHASH] = (*StackMachine).opTodo
	complexFuncs[SELFBALANCE] = (*StackMachine).opTodo

	complexFuncs[MSTORE8] = (*StackMachine).opTodo
	complexFuncs[SLOAD] = (*StackMachine).opTodo
	complexFuncs[SSTORE] = (*StackMachine).opTodo
	complexFuncs[JUMP] = (*StackMachine).opJump
	complexFuncs[JUMPI] = (*StackMachine).opJumpi
	complexFuncs[PC] = (*StackMachine).opTodo

	complexFuncs[GAS] = (*StackMachine).opTodo
	complexFuncs[SELFBALANCE] = (*StackMachine).opTodo

	complexFuncs[LOG0] = (*StackMachine).opTodo
	complexFuncs[LOG1] = (*StackMachine).opTodo
	complexFuncs[LOG2] = (*StackMachine).opTodo
	complexFuncs[LOG3] = (*StackMachine).opTodo
	complexFuncs[LOG4] = (*StackMachine).opTodo

	complexFuncs[CREATE] = (*StackMachine).opTodo
	complexFuncs[CALL] = (*StackMachine).opTodo
	complexFuncs[RETURN] = (*StackMachine).opTodo
	complexFuncs[DELEGATECALL] = (*StackMachine).opTodo
	complexFuncs[CREATE2] = (*StackMachine).opTodo

	complexFuncs[STATICCALL] = (*StackMachine).opTodo
	complexFuncs[REVERT] = (*StackMachine).opTodo
	complexFuncs[SELFDESTRUCT] = (*StackMachine).opTodo

	staticGasCosts[PUSH3] = 5
	staticGasCosts[JUMPDEST] = 0
	staticGasCosts[PUSH1] = 5
	staticGasCosts[SWAP1] = 5
	staticGasCosts[SUB] = 5
	staticGasCosts[DUP1] = 5
	staticGasCosts[JUMPI] = 5
	staticGasCosts[STOP] = 5

}

func (machine *StackMachine) DispatchSimple(op OpCode, pc *uint64) (valid bool) {
	if fn := simpleFuncs[op]; fn != nil {
		fn(machine)
		return true
	}
	if op >= PUSH1 && op <= PUSH32 {
		start := int(*pc) + 1
		code := machine.callCtx.Code
		l := int(1 + byte(op-PUSH1))
		end := start + l
		if codeLen := len(code); codeLen < end {
			end = int(codeLen)
		}
		machine.PushBytes(code[start:end])
		*pc += uint64(l)
		return true
	} else if op >= DUP1 && op <= DUP16 {
		machine.Dup(int(op - DUP1))
		return true
	} else if op >= SWAP1 && op <= SWAP16 {
		machine.Swap(int(1+op - SWAP1))
		return true
	}
	return false
}

// Call is the entry-point for execution, and is used for the initial outer call
// as well as any internal sub-calls, such as DELEGATECALL and STATICCALL
func (machine *StackMachine) Call(code []byte, address, caller [20]byte, value *uint256.Int, readOnly bool, gas uint64) ([]byte, error) {
	if len(code) == 0 {
		return nil, nil
	}
	// New context
	machine.NewContext(code, value)
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
		abort        bool
		op           OpCode
		//steps        uint64
	)

	for ; !abort; pc++ {
		//steps++
		//if steps == 6 + 1000000{
		//	break
		//}
		op = OpCode(code[pc])
		//fmt.Printf("pc %d op %v\n", pc, op.String())
		if op == 0 {
			break
		}
		// An opcode with a negative static gas value is not valid
		if cost := machine.staticGasCost[op]; cost < 0 {
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
		} else {
			opStaticCost = cost
		}
		// Validate stack
		if sLen := machine.StackDepth(); sLen < MinStack[op] {
			return nil, fmt.Errorf("stack underflow (%d <=> %d)", sLen, MinStack[op])
		} else if sLen > MaxStack[op] {
			return nil, fmt.Errorf("stack limit reached %d (%d)", sLen, MaxStack[op])
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

		if machine.DispatchSimple(op, &pc) {
			continue
		}
		// All 'simple' opcodes are done. Remaining are calls, tx-context stackbuffer,
		// and a few others.
		// The non-simple opcodes may
		// -- require state access,
		// -- require access to gas,
		// -- have dynamic costs
		// -- return errors
		var err error
		// Complex ops remaining. Complex ops may return error
		if fn := complexFuncs[op]; fn != nil {
			gas, err = fn(machine, &pc, gas)
		}
		if err != nil {
			// TODO handle this?
			return nil, err
		}
	}
	return nil, nil
}

func calcMemSize64(a, b *uint256.Int) (uint64, bool) {
	panic("implement me")
}

func memoryGasCost(a, b *uint256.Int) (uint64, bool) {
	panic("implement me")
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

//func (machine *StackMachine) dynamicGasCall(value *uint256.Int, destination [20]byte) {
//	var (
//		gasCost        uint64
//		transfersValue = !value.IsZero()
//	)
//
//	// Determine how much gas is needed for the memory expansion (if any)
//	memoryGas, err := memoryGasCost(machine.m, memorySize)
//	if err != nil {
//		return 0, 0, err
//	}
//	// Determine extra gas that needs to be paid for creating a new account
//	if isEip158 {
//		if transfersValue && machine.StateDB.Empty(destination) {
//			gasCost += params.CallNewAccountGas
//		}
//	} else if !machine.StateDB.Exist(destination) {
//		gasCost += params.CallNewAccountGas
//	}
//	// Extra fee for transferring ether
//	if transfersValue {
//		gasCost += params.CallValueTransferGas
//	}
//	// Add the two factors
//	var overflow bool
//	if gasCost, overflow = math.SafeAdd(gasCost, memoryGas); overflow {
//		return 0, 0, errGasUintOverflow
//	}
//	// Apply the 63/64:ths rule, to determine how much gas is actually
//	// passed along to the callee
//	sentGas, err = callGas(isEip150, availableGas, gasCost, desiredGas)
//	if err != nil {
//		return 0, 0, err
//	}
//	// Check if we have sufficient gas
//	if gasCost, overflow = math.SafeAdd(availableGas, sentGas); overflow {
//		return 0, 0, errGasUintOverflow
//	}
//	return gasCost, sentGas, nil
//}
//
//func (machine *StackMachine) opCall(availableGas uint64) (remainingGas uint64, err error) {
//
//	var (
//		desiredGas  = uint256.NewInt()
//		destination [20]byte
//		value       = uint256.NewInt()
//		inOffset    = uint256.NewInt()
//		inSize      = uint256.NewInt()
//		retOffset   = uint256.NewInt()
//		retSize     = uint256.NewInt()
//	)
//	machine.PopUint(desiredGas)
//	machine.PopBytes20(destination)
//	machine.PopUint(value)
//	machine.PopUints(inOffset, inSize)
//	machine.PopUints(retOffset, retSize)
//
//	// We do the memory expansion calculation and dynamic gas calc here
//	memSize, overflow := machine.memoryCall(inOffset, inSize, retOffset, retSize)
//	if overflow {
//		return 0, errGasUintOverflow
//	}
//	dynamicGasCost, calleeGas, err := machine.dynamicGasCall(value, destination, memSize)
//	if err != nil {
//		return 0, err
//	}
//	if availableGas < dynamicGasCost {
//		return 0, ErrOutOfGas
//	}
//	availableGas -= dynamicGasCost
//	// At this point, we know there is sufficient gas, so it's time to
//	// expand the memory
//	if memSize > 0 {
//		machine.callCtx.memory.Resize(memorySize)
//	}
//	// Execute the CALL
//	ret, returnedGas, err := machine.Call(code, address, caller, value, false, calleeGas)
//	if err != nil {
//		machine.PushZero()
//	} else {
//		machine.PushBytes(oneByte32)
//	}
//	if err == nil || err == errExecutionReverted {
//		machine.callCtx.memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
//		machine.returnData = ret
//	}
//	return availableGas + returnedGas, nil
//}

func (machine *StackMachine) opTodo(pc *uint64, gas uint64) (uint64, error) {
	panic("not implemented")
}
