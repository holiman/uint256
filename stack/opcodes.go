// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package stack

import (
	"fmt"
)

// OpCode is an EVM opcode
type OpCode byte

const (
	// 0x0 range - arithmetic ops.
	STOP       = OpCode(0x00)
	ADD        = OpCode(0x01)
	MUL        = OpCode(0x02)
	SUB        = OpCode(0x03)
	DIV        = OpCode(0x04)
	SDIV       = OpCode(0x05)
	MOD        = OpCode(0x06)
	SMOD       = OpCode(0x07)
	ADDMOD     = OpCode(0x08)
	MULMOD     = OpCode(0x09)
	EXP        = OpCode(0x0a)
	SIGNEXTEND = OpCode(0x0b)

	// 0x10 range - comparison ops.
	LT     = OpCode(0x10)
	GT     = OpCode(0x11)
	SLT    = OpCode(0x12)
	SGT    = OpCode(0x13)
	EQ     = OpCode(0x14)
	ISZERO = OpCode(0x15)
	AND    = OpCode(0x16)
	OR     = OpCode(0x17)
	XOR    = OpCode(0x18)
	NOT    = OpCode(0x19)
	BYTE   = OpCode(0x1a)
	SHL    = OpCode(0x1b)
	SHR    = OpCode(0x1c)
	SAR    = OpCode(0x1d)

	// 0x20 range - sha3.
	SHA3 = OpCode(0x20)

	// 0x30 range - closure state.
	ADDRESS        = OpCode(0x30)
	BALANCE        = OpCode(0x31)
	ORIGIN         = OpCode(0x32)
	CALLER         = OpCode(0x33)
	CALLVALUE      = OpCode(0x34)
	CALLDATALOAD   = OpCode(0x35)
	CALLDATASIZE   = OpCode(0x36)
	CALLDATACOPY   = OpCode(0x37)
	CODESIZE       = OpCode(0x38)
	CODECOPY       = OpCode(0x39)
	GASPRICE       = OpCode(0x3a)
	EXTCODESIZE    = OpCode(0x3b)
	EXTCODECOPY    = OpCode(0x3c)
	RETURNDATASIZE = OpCode(0x3d)
	RETURNDATACOPY = OpCode(0x3e)
	EXTCODEHASH    = OpCode(0x3f)

	// 0x40 range - block operations + selfbalance and chainid
	BLOCKHASH   = OpCode(0x40)
	COINBASE    = OpCode(0x41)
	TIMESTAMP   = OpCode(0x42)
	NUMBER      = OpCode(0x43)
	DIFFICULTY  = OpCode(0x44)
	GASLIMIT    = OpCode(0x45)
	CHAINID     = OpCode(0x46)
	SELFBALANCE = OpCode(0x47)

	// 0x50 range - 'storage' and execution.
	POP      = OpCode(0x50)
	MLOAD    = OpCode(0x51)
	MSTORE   = OpCode(0x52)
	MSTORE8  = OpCode(0x53)
	SLOAD    = OpCode(0x54)
	SSTORE   = OpCode(0x55)
	JUMP     = OpCode(0x56)
	JUMPI    = OpCode(0x57)
	PC       = OpCode(0x58)
	MSIZE    = OpCode(0x59)
	GAS      = OpCode(0x5a)
	JUMPDEST = OpCode(0x5b)

	// 0x60 range.
	PUSH1  = OpCode(0x60)
	PUSH2  = OpCode(0x61)
	PUSH3  = OpCode(0x62)
	PUSH4  = OpCode(0x63)
	PUSH5  = OpCode(0x64)
	PUSH6  = OpCode(0x65)
	PUSH7  = OpCode(0x66)
	PUSH8  = OpCode(0x67)
	PUSH9  = OpCode(0x68)
	PUSH10 = OpCode(0x69)
	PUSH11 = OpCode(0x6a)
	PUSH12 = OpCode(0x6b)
	PUSH13 = OpCode(0x6c)
	PUSH14 = OpCode(0x6d)
	PUSH15 = OpCode(0x6e)
	PUSH16 = OpCode(0x6f)
	PUSH17 = OpCode(0x70)
	PUSH18 = OpCode(0x71)
	PUSH19 = OpCode(0x72)
	PUSH20 = OpCode(0x73)
	PUSH21 = OpCode(0x74)
	PUSH22 = OpCode(0x75)
	PUSH23 = OpCode(0x76)
	PUSH24 = OpCode(0x77)
	PUSH25 = OpCode(0x78)
	PUSH26 = OpCode(0x79)
	PUSH27 = OpCode(0x7a)
	PUSH28 = OpCode(0x7b)
	PUSH29 = OpCode(0x7c)
	PUSH30 = OpCode(0x7d)
	PUSH31 = OpCode(0x7e)
	PUSH32 = OpCode(0x7f)
	DUP1   = OpCode(0x80)
	DUP2   = OpCode(0x81)
	DUP3   = OpCode(0x82)
	DUP4   = OpCode(0x83)
	DUP5   = OpCode(0x84)
	DUP6   = OpCode(0x85)
	DUP7   = OpCode(0x86)
	DUP8   = OpCode(0x87)
	DUP9   = OpCode(0x88)
	DUP10  = OpCode(0x89)
	DUP11  = OpCode(0x8a)
	DUP12  = OpCode(0x8b)
	DUP13  = OpCode(0x8c)
	DUP14  = OpCode(0x8d)
	DUP15  = OpCode(0x8e)
	DUP16  = OpCode(0x8f)
	SWAP1  = OpCode(0x90)
	SWAP2  = OpCode(0x91)
	SWAP3  = OpCode(0x92)
	SWAP4  = OpCode(0x93)
	SWAP5  = OpCode(0x94)
	SWAP6  = OpCode(0x95)
	SWAP7  = OpCode(0x96)
	SWAP8  = OpCode(0x97)
	SWAP9  = OpCode(0x98)
	SWAP10 = OpCode(0x99)
	SWAP11 = OpCode(0x9a)
	SWAP12 = OpCode(0x9b)
	SWAP13 = OpCode(0x9c)
	SWAP14 = OpCode(0x9d)
	SWAP15 = OpCode(0x9e)
	SWAP16 = OpCode(0x9f)

	// 0xa0 range - logging ops.
	LOG0 = OpCode(0xa0)
	LOG1 = OpCode(0xa1)
	LOG2 = OpCode(0xa2)
	LOG3 = OpCode(0xa3)
	LOG4 = OpCode(0xa4)

	// 0xf0 range - closures.
	CREATE       = OpCode(0xf0)
	CALL         = OpCode(0xf1)
	CALLCODE     = OpCode(0xf2)
	RETURN       = OpCode(0xf3)
	DELEGATECALL = OpCode(0xf4)
	CREATE2      = OpCode(0xf5)

	STATICCALL   = OpCode(0xfa)
	REVERT       = OpCode(0xfd)
	SELFDESTRUCT = OpCode(0xff)
)

// Since the opcodes aren't all in order we can't use a regular slice.
var opCodeToString = map[OpCode]string{
	// 0x0 range - arithmetic ops.
	STOP:       "STOP",
	ADD:        "ADD",
	MUL:        "MUL",
	SUB:        "SUB",
	DIV:        "DIV",
	SDIV:       "SDIV",
	MOD:        "MOD",
	SMOD:       "SMOD",
	EXP:        "EXP",
	NOT:        "NOT",
	LT:         "LT",
	GT:         "GT",
	SLT:        "SLT",
	SGT:        "SGT",
	EQ:         "EQ",
	ISZERO:     "ISZERO",
	SIGNEXTEND: "SIGNEXTEND",

	// 0x10 range - bit ops.
	AND:    "AND",
	OR:     "OR",
	XOR:    "XOR",
	BYTE:   "BYTE",
	SHL:    "SHL",
	SHR:    "SHR",
	SAR:    "SAR",
	ADDMOD: "ADDMOD",
	MULMOD: "MULMOD",

	// 0x20 range - crypto.
	SHA3: "SHA3",

	// 0x30 range - closure state.
	ADDRESS:        "ADDRESS",
	BALANCE:        "BALANCE",
	ORIGIN:         "ORIGIN",
	CALLER:         "CALLER",
	CALLVALUE:      "CALLVALUE",
	CALLDATALOAD:   "CALLDATALOAD",
	CALLDATASIZE:   "CALLDATASIZE",
	CALLDATACOPY:   "CALLDATACOPY",
	CODESIZE:       "CODESIZE",
	CODECOPY:       "CODECOPY",
	GASPRICE:       "GASPRICE",
	EXTCODESIZE:    "EXTCODESIZE",
	EXTCODECOPY:    "EXTCODECOPY",
	RETURNDATASIZE: "RETURNDATASIZE",
	RETURNDATACOPY: "RETURNDATACOPY",
	EXTCODEHASH:    "EXTCODEHASH",

	// 0x40 range - block operations.
	BLOCKHASH:   "BLOCKHASH",
	COINBASE:    "COINBASE",
	TIMESTAMP:   "TIMESTAMP",
	NUMBER:      "NUMBER",
	DIFFICULTY:  "DIFFICULTY",
	GASLIMIT:    "GASLIMIT",
	CHAINID:     "CHAINID",
	SELFBALANCE: "SELFBALANCE",

	// 0x50 range - 'storage' and execution.
	POP: "POP",
	//DUP:     "DUP",
	//SWAP:    "SWAP",
	MLOAD:    "MLOAD",
	MSTORE:   "MSTORE",
	MSTORE8:  "MSTORE8",
	SLOAD:    "SLOAD",
	SSTORE:   "SSTORE",
	JUMP:     "JUMP",
	JUMPI:    "JUMPI",
	PC:       "PC",
	MSIZE:    "MSIZE",
	GAS:      "GAS",
	JUMPDEST: "JUMPDEST",

	// 0x60 range - push.
	PUSH1:  "PUSH1",
	PUSH2:  "PUSH2",
	PUSH3:  "PUSH3",
	PUSH4:  "PUSH4",
	PUSH5:  "PUSH5",
	PUSH6:  "PUSH6",
	PUSH7:  "PUSH7",
	PUSH8:  "PUSH8",
	PUSH9:  "PUSH9",
	PUSH10: "PUSH10",
	PUSH11: "PUSH11",
	PUSH12: "PUSH12",
	PUSH13: "PUSH13",
	PUSH14: "PUSH14",
	PUSH15: "PUSH15",
	PUSH16: "PUSH16",
	PUSH17: "PUSH17",
	PUSH18: "PUSH18",
	PUSH19: "PUSH19",
	PUSH20: "PUSH20",
	PUSH21: "PUSH21",
	PUSH22: "PUSH22",
	PUSH23: "PUSH23",
	PUSH24: "PUSH24",
	PUSH25: "PUSH25",
	PUSH26: "PUSH26",
	PUSH27: "PUSH27",
	PUSH28: "PUSH28",
	PUSH29: "PUSH29",
	PUSH30: "PUSH30",
	PUSH31: "PUSH31",
	PUSH32: "PUSH32",

	DUP1:  "DUP1",
	DUP2:  "DUP2",
	DUP3:  "DUP3",
	DUP4:  "DUP4",
	DUP5:  "DUP5",
	DUP6:  "DUP6",
	DUP7:  "DUP7",
	DUP8:  "DUP8",
	DUP9:  "DUP9",
	DUP10: "DUP10",
	DUP11: "DUP11",
	DUP12: "DUP12",
	DUP13: "DUP13",
	DUP14: "DUP14",
	DUP15: "DUP15",
	DUP16: "DUP16",

	SWAP1:  "SWAP1",
	SWAP2:  "SWAP2",
	SWAP3:  "SWAP3",
	SWAP4:  "SWAP4",
	SWAP5:  "SWAP5",
	SWAP6:  "SWAP6",
	SWAP7:  "SWAP7",
	SWAP8:  "SWAP8",
	SWAP9:  "SWAP9",
	SWAP10: "SWAP10",
	SWAP11: "SWAP11",
	SWAP12: "SWAP12",
	SWAP13: "SWAP13",
	SWAP14: "SWAP14",
	SWAP15: "SWAP15",
	SWAP16: "SWAP16",
	LOG0:   "LOG0",
	LOG1:   "LOG1",
	LOG2:   "LOG2",
	LOG3:   "LOG3",
	LOG4:   "LOG4",

	// 0xf0 range.
	CREATE:       "CREATE",
	CALL:         "CALL",
	RETURN:       "RETURN",
	CALLCODE:     "CALLCODE",
	DELEGATECALL: "DELEGATECALL",
	CREATE2:      "CREATE2",
	STATICCALL:   "STATICCALL",
	REVERT:       "REVERT",
	SELFDESTRUCT: "SELFDESTRUCT",

}

func (op OpCode) String() string {
	str := opCodeToString[op]
	if len(str) == 0 {
		return fmt.Sprintf("Missing opcode 0x%x", int(op))
	}

	return str
}

var stringToOp = map[string]OpCode{
	"STOP":           STOP,
	"ADD":            ADD,
	"MUL":            MUL,
	"SUB":            SUB,
	"DIV":            DIV,
	"SDIV":           SDIV,
	"MOD":            MOD,
	"SMOD":           SMOD,
	"EXP":            EXP,
	"NOT":            NOT,
	"LT":             LT,
	"GT":             GT,
	"SLT":            SLT,
	"SGT":            SGT,
	"EQ":             EQ,
	"ISZERO":         ISZERO,
	"SIGNEXTEND":     SIGNEXTEND,
	"AND":            AND,
	"OR":             OR,
	"XOR":            XOR,
	"BYTE":           BYTE,
	"SHL":            SHL,
	"SHR":            SHR,
	"SAR":            SAR,
	"ADDMOD":         ADDMOD,
	"MULMOD":         MULMOD,
	"SHA3":           SHA3,
	"ADDRESS":        ADDRESS,
	"BALANCE":        BALANCE,
	"ORIGIN":         ORIGIN,
	"CALLER":         CALLER,
	"CALLVALUE":      CALLVALUE,
	"CALLDATALOAD":   CALLDATALOAD,
	"CALLDATASIZE":   CALLDATASIZE,
	"CALLDATACOPY":   CALLDATACOPY,
	"CHAINID":        CHAINID,
	"DELEGATECALL":   DELEGATECALL,
	"STATICCALL":     STATICCALL,
	"CODESIZE":       CODESIZE,
	"CODECOPY":       CODECOPY,
	"GASPRICE":       GASPRICE,
	"EXTCODESIZE":    EXTCODESIZE,
	"EXTCODECOPY":    EXTCODECOPY,
	"RETURNDATASIZE": RETURNDATASIZE,
	"RETURNDATACOPY": RETURNDATACOPY,
	"EXTCODEHASH":    EXTCODEHASH,
	"BLOCKHASH":      BLOCKHASH,
	"COINBASE":       COINBASE,
	"TIMESTAMP":      TIMESTAMP,
	"NUMBER":         NUMBER,
	"DIFFICULTY":     DIFFICULTY,
	"GASLIMIT":       GASLIMIT,
	"SELFBALANCE":    SELFBALANCE,
	"POP":            POP,
	"MLOAD":          MLOAD,
	"MSTORE":         MSTORE,
	"MSTORE8":        MSTORE8,
	"SLOAD":          SLOAD,
	"SSTORE":         SSTORE,
	"JUMP":           JUMP,
	"JUMPI":          JUMPI,
	"PC":             PC,
	"MSIZE":          MSIZE,
	"GAS":            GAS,
	"JUMPDEST":       JUMPDEST,
	"PUSH1":          PUSH1,
	"PUSH2":          PUSH2,
	"PUSH3":          PUSH3,
	"PUSH4":          PUSH4,
	"PUSH5":          PUSH5,
	"PUSH6":          PUSH6,
	"PUSH7":          PUSH7,
	"PUSH8":          PUSH8,
	"PUSH9":          PUSH9,
	"PUSH10":         PUSH10,
	"PUSH11":         PUSH11,
	"PUSH12":         PUSH12,
	"PUSH13":         PUSH13,
	"PUSH14":         PUSH14,
	"PUSH15":         PUSH15,
	"PUSH16":         PUSH16,
	"PUSH17":         PUSH17,
	"PUSH18":         PUSH18,
	"PUSH19":         PUSH19,
	"PUSH20":         PUSH20,
	"PUSH21":         PUSH21,
	"PUSH22":         PUSH22,
	"PUSH23":         PUSH23,
	"PUSH24":         PUSH24,
	"PUSH25":         PUSH25,
	"PUSH26":         PUSH26,
	"PUSH27":         PUSH27,
	"PUSH28":         PUSH28,
	"PUSH29":         PUSH29,
	"PUSH30":         PUSH30,
	"PUSH31":         PUSH31,
	"PUSH32":         PUSH32,
	"DUP1":           DUP1,
	"DUP2":           DUP2,
	"DUP3":           DUP3,
	"DUP4":           DUP4,
	"DUP5":           DUP5,
	"DUP6":           DUP6,
	"DUP7":           DUP7,
	"DUP8":           DUP8,
	"DUP9":           DUP9,
	"DUP10":          DUP10,
	"DUP11":          DUP11,
	"DUP12":          DUP12,
	"DUP13":          DUP13,
	"DUP14":          DUP14,
	"DUP15":          DUP15,
	"DUP16":          DUP16,
	"SWAP1":          SWAP1,
	"SWAP2":          SWAP2,
	"SWAP3":          SWAP3,
	"SWAP4":          SWAP4,
	"SWAP5":          SWAP5,
	"SWAP6":          SWAP6,
	"SWAP7":          SWAP7,
	"SWAP8":          SWAP8,
	"SWAP9":          SWAP9,
	"SWAP10":         SWAP10,
	"SWAP11":         SWAP11,
	"SWAP12":         SWAP12,
	"SWAP13":         SWAP13,
	"SWAP14":         SWAP14,
	"SWAP15":         SWAP15,
	"SWAP16":         SWAP16,
	"LOG0":           LOG0,
	"LOG1":           LOG1,
	"LOG2":           LOG2,
	"LOG3":           LOG3,
	"LOG4":           LOG4,
	"CREATE":         CREATE,
	"CREATE2":        CREATE2,
	"CALL":           CALL,
	"RETURN":         RETURN,
	"CALLCODE":       CALLCODE,
	"REVERT":         REVERT,
	"SELFDESTRUCT":   SELFDESTRUCT,
}

// StringToOp finds the opcode whose name is stored in `str`.
func StringToOp(str string) OpCode {
	return stringToOp[str]
}

// TODO initialize these
var MinStack [256]int
var MaxStack [256]int
func init(){
	for i := 0; i < 256; i++{
		MaxStack[i] = 10
	}
}