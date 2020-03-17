package stack

import "math/big"

// The ROM implements a read-only memory, meant to be used for block-constants
// and tx-constants
// This implementation is inspired by
// https://corepaper.org/ethereum/evm/#evm-rom
// Which defines the following
//
// ADDRESS (0x30): READROM 0x0 Push index 0 of read-only memory onto stack.
// ORIGIN (0x32): READROM 0x1 Push index 1 of read-only memory onto stack.
// CALLER (0x33): READROM 0x3 Push index 2 of read-only memory onto stack.
// CALLVALUE (0x34): READROM 0x4 Push index 3 of read-only memory onto stack.
// GASPRICE (0x3a): READROM 0x5 Push index 4 of read-only memory onto stack.
// COINBASE (0x41): READROM 0x6 Push index 5 of read-only memory onto stack.
// TIMESTAMP (0x42): READROM 0x7 Push index 6 of read-only memory onto stack.
// NUMBER (0x43): READROM 0x8 Push index 7 of read-only memory onto stack.
// DIFFICULTY (0x44): READROM 0x9 Push index 8 of read-only memory onto stack.
// GASLIMIT (0x45): READROM 0xa Push index 9 of read-only memory onto stack.
// CHAINID (0x46): READROM 0xb Push index 10 of read-only memory onto stack.
// SELFBALANCE (0x47): READROM 0xc Push index 11 of read-only memory onto stack.
//
// However, this implementation uses only the ones that are defined per-block:

// ORIGIN (0x32): READROM 0x1 Push index 1 of read-only memory onto stack.
// GASPRICE (0x3a): READROM 0x5 Push index 4 of read-only memory onto stack.
// COINBASE (0x41): READROM 0x6 Push index 5 of read-only memory onto stack.
// TIMESTAMP (0x42): READROM 0x7 Push index 6 of read-only memory onto stack.
// NUMBER (0x43): READROM 0x8 Push index 7 of read-only memory onto stack.
// DIFFICULTY (0x44): READROM 0x9 Push index 8 of read-only memory onto stack.
// CHAINID (0x46): READROM 0xb Push index 10 of read-only memory onto stack.
//
// Remaining ones are implemented elsewhere:
//
// ADDRESS (0x30): READROM 0x0 Push index 0 of read-only memory onto stack.
// CALLER (0x33): READROM 0x3 Push index 2 of read-only memory onto stack.
// CALLVALUE (0x34): READROM 0x4 Push index 3 of read-only memory onto stack.
// SELFBALANCE (0x47): READROM 0xc Push index 11 of read-only memory onto stack.

const (
	romOrigin uint = 32 * iota
	romGasPrice
	romCoinbase
	romTimestamp
	romNumber
	romDifficulty
	romChainId

	romUnused // EOF
)

type readonlyMem []byte

// Create readonly-mem containing per-block constants
func newReadonlyMem() readonlyMem {
	rom := make(readonlyMem, romUnused)
	return rom
}

func (rom readonlyMem) setBlockConstants(coinbase, timestamp, number,
	difficulty, chainid *big.Int) {
	copy(rom[romCoinbase:romCoinbase+32], coinbase.Bytes())
	copy(rom[romTimestamp:romTimestamp+32], timestamp.Bytes())
	copy(rom[romNumber:romNumber+32], number.Bytes())
	copy(rom[romDifficulty:romDifficulty+32], difficulty.Bytes())
	copy(rom[romChainId:romChainId+32], chainid.Bytes())

}

// setTxConstants initializes the readonly-memory with transaction-specific
// stackbuffer
func (rom readonlyMem) setTxConstants(origin, gasPrice *uint256.Int) readonlyMem {
	copy(rom[romOrigin:romOrigin+32], origin.Bytes())
	copy(rom[romGasPrice:romGasPrice+32], gasPrice.Bytes())
	return rom
}

func (machine *StackMachine) opCoinbase() {
	machine.PushBytes(machine.rom[romCoinbase : romCoinbase+32])
}

func (machine *StackMachine) opTimestamp() {
	machine.PushBytes(machine.rom[romTimestamp : romTimestamp+32])
}

func (machine *StackMachine) opNumber() {
	machine.PushBytes(machine.rom[romNumber : romNumber+32])
}
func (machine *StackMachine) opDifficulty() {
	machine.PushBytes(machine.rom[romDifficulty : romDifficulty+32])
}
func (machine *StackMachine) opChainId() {
	machine.PushBytes(machine.rom[romChainId : romChainId+32])
}

func (machine *StackMachine) opOrigin() {
	machine.PushBytes(machine.rom[romOrigin : romOrigin+32])
}

func (machine *StackMachine) opGasprice() {
	machine.PushBytes(machine.rom[romGasPrice : romGasPrice+32])
}
