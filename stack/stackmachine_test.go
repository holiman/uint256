package stack

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/holiman/uint256"
)

func toHex(h string) []byte {
	b, _ := hex.DecodeString(h)
	return b
}

func TestStackBasics(t *testing.T) {
	machine := New()
	machine.NewContext(nil,nil)
	{
		machine.PushBytes(toHex("deadbeef"))
		machine.NewContext(nil,nil)
		{
			machine.PushBytes(toHex("deadbeef"))
			machine.PushBytes(toHex("deadbeef"))
			if got, exp := machine.callCtx.stackHead,3 * 32; got != exp {
				t.Fatalf("machine head wrong, got %v exp %v", got, exp)
			}
			machine.Pop()
			if got, exp := machine.callCtx.stackHead,2 * 32; got != exp {
				t.Fatalf("machine head wrong, got %v exp %v", got, exp)
			}
			machine.PushBytes(toHex("deadbeef0000000000000000"))
			machine.PushBytes(toHex("ffaa112233deadbeef000cafebabe00000deadbeef00000cafebabe000deadbeef01020304"))
			d := machine.PopBytes32([32]byte{})
			if exp := toHex("deadbeef000cafebabe00000deadbeef00000cafebabe000deadbeef01020304"); !bytes.Equal(d[:],exp){
				t.Fatalf("err, got %x exp %x", d, exp)
			}
			machine.DropContext()
		}
		machine.DropContext()
	}
}

func TestStackInts(t *testing.T) {
	stack := New()
	stack.NewContext(nil, nil)
	{
		stack.PushBytes(toHex("01"))
		stack.PushBytes(toHex("02"))
		a := uint256.NewInt()
		b := uint256.NewInt()
		stack.PopUint(a)
		stack.PopUint(b)
		a.Add(a, b)
		if v, oflow := a.Uint64WithOverflow(); oflow || v != 3 {
			t.Fatalf("got %v %v, exp %v %v", v, oflow, 3, false)
		}
		stack.PushUint(&uint256.Int{1, 1, 1, 1})
		stack.Pop()
		stack.PushUint(&uint256.Int{})
		stack.PopUint(a)
		if v, oflow := a.Uint64WithOverflow(); oflow || v != 0 {
			t.Fatalf("got %v %v, exp %v %v", v, oflow, 0, false)
		}
		stack.PushUint64(65)
		stack.PopUint(a)
		if v, oflow := a.Uint64WithOverflow(); oflow || v != 65 {
			t.Fatalf("got %v %v, exp %v %v", v, oflow, 65, false)
		}
		stack.PushBytes(toHex("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
		stack.PushBytes(toHex("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))
		stack.PushBytes(toHex("cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"))
		stack.PushBytes(toHex("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"))
		stack.Swap(1) // should be [a,b, d, c]
		stack.Swap(2) // should be [a,c, d, b]
		stack.Swap(3) // should be [b, c, d, a]
		stack.PrettyPrint()
		stack.Dup(2) // should be [b, c, d, a, c]
		stack.PrettyPrint()
		stack.DropContext()
	}
}

