package stack

func (machine *StackMachine) opAdd() {
	machine.PopUints(machine.x, machine.y)
	machine.z.Add(machine.x, machine.y)
	machine.PushUint(machine.z)
}

func (machine *StackMachine) opSub() {
	machine.PopUints(machine.x, machine.y)
	machine.z.Sub(machine.x, machine.y)
	machine.PushUint(machine.z)
}

func (machine *StackMachine) opMul() {
	machine.PopUints(machine.x, machine.y)
	machine.z.Mul(machine.x, machine.y)
	machine.PushUint(machine.z)
}

func (machine *StackMachine) opDiv() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Div(machine.x, machine.y))
}

func (machine *StackMachine) opSdiv() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Sdiv(machine.x, machine.y))
}

func (machine *StackMachine) opMod() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Mod(machine.x, machine.y))
}

func (machine *StackMachine) opSMod() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Smod(machine.x, machine.y))
}

func (machine *StackMachine) opExp() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Exp(machine.x, machine.y))
}

func (machine *StackMachine) opSignExtend() {
	machine.PopUints(machine.x, machine.y)
	machine.z.SignExtend(machine.x, machine.y)
	machine.PushUint(machine.z)
}

func (machine *StackMachine) opNot() {
	machine.PopUint(machine.x)
	machine.PushUint(machine.x.Not())
}

func (machine *StackMachine) opLt() {
	machine.PopUints(machine.x, machine.y)
	machine.PushBool(machine.x.Lt(machine.y))
}

func (machine *StackMachine) opGt() {
	machine.PopUints(machine.x, machine.y)
	machine.PushBool(machine.x.Gt(machine.y))
}

func (machine *StackMachine) opSlt() {
	machine.PopUints(machine.x, machine.y)
	machine.PushBool(machine.x.Slt(machine.y))
}

func (machine *StackMachine) opSgt() {
	machine.PopUints(machine.x, machine.y)
	machine.PushBool(machine.x.Sgt(machine.y))
}

func (machine *StackMachine) opEq() {
	machine.PopUints(machine.x, machine.y)
	machine.PushBool(machine.x.Eq(machine.y))
}

func (machine *StackMachine) opIszero() {
	machine.PopUint(machine.x)
	machine.PushBool(machine.x.IsZero())
}

func (machine *StackMachine) opAnd() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.And(machine.x, machine.y))
}

func (machine *StackMachine) opOr() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Or(machine.x, machine.y))
}

func (machine *StackMachine) opXor() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.z.Xor(machine.x, machine.y))
}

func (machine *StackMachine) opByte() {
	machine.PopUints(machine.x, machine.y)
	machine.PushUint(machine.y.Byte(machine.x))
}

func (machine *StackMachine) opAddmod() {
	machine.PopUints(machine.x, machine.y)
	machine.PopUint(machine.z)
	if machine.z.IsZero() {
		machine.PushZero()
		return
	}
	machine.z.AddMod(machine.x, machine.y, machine.z)
	machine.PushUint(machine.z)
}

func (machine *StackMachine) opMulmod() {
	machine.PopUints(machine.x, machine.y)
	machine.PopUint(machine.z)
	machine.z.MulMod(machine.x, machine.y, machine.z)
	machine.PushUint(machine.z)
}

func (machine *StackMachine) opSHL() {
	machine.PopUints(machine.x, machine.y)
	if machine.x.LtUint64(256) {
		machine.PushUint(machine.y.Lsh(machine.y, uint(machine.x.Uint64())))
		return
	}
	machine.PushZero()
}

func (machine *StackMachine) opSHR() {
	machine.PopUints(machine.x, machine.y)
	if machine.x.LtUint64(256) {
		machine.PushUint(machine.y.Rsh(machine.y, uint(machine.x.Uint64())))
		return
	}
	machine.PushZero()
}

func (machine *StackMachine) opSAR() {
	machine.PopUints(machine.x, machine.y)
	if machine.x.GtUint64(256) {
		if machine.y.Sign() >= 0 {
			machine.PushZero()
			return
		}
		machine.y.SetAllOne()
		machine.PushUint(machine.y)
		return
	}
	n := uint(machine.x.Uint64())
	machine.y.Srsh(machine.y, n)
	machine.PushUint(machine.y)
}

func (machine *StackMachine) opMload() {
	machine.PopUint(machine.x)
	offset := machine.x.Int64()
	available := int64(len(machine.callCtx.memory))
	if available < offset {
		machine.PushZero()
	} else if available < offset+32 {
		machine.PushBytes(machine.callCtx.memory[offset:available])
	} else {
		machine.PushBytes(machine.callCtx.memory[offset : offset+32])
	}
}
func (machine *StackMachine) opMstore() {
	// memStart , value
	machine.PopUints(machine.x, machine.y)
	offset := machine.x.Int64()
	// This will panic if memory is not already expanded
	machine.y.WriteToSlice(machine.callCtx.memory[offset:])
}

func (machine *StackMachine) opMstore8() {
	// memStart , value
	machine.PopUints(machine.x, machine.y)
	offset := machine.x.Int64()
	// This will panic if memory is not already expanded
	machine.callCtx.memory[offset] = byte(machine.y.Int64() & 0xFF)
}

func (machine *StackMachine) opMsize() {
	machine.PushUint64(uint64(len(machine.callCtx.memory)))
}
