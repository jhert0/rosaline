package chess

import (
	"testing"
)

func TestSetBit(t *testing.T) {
	bb := NewBitBoard(0)
	bb.SetBit(1)
	if bb != 2 {
		t.Fatalf("%v: expected: 2 got %v", t.Name(), bb)
	}
}

func TestClearBit(t *testing.T) {
	bb := NewBitBoard(1)
	bb.ClearBit(0)
	if bb != 0 {
		t.Fatalf("%v: expected: 0 got %v", t.Name(), bb)
	}
}

func TestBitSet(t *testing.T) {
	bb := NewBitBoard(1)
	if bb.BitSet(1) {
		t.Fatalf("%v: expected: true got %v", t.Name(), bb)
	}
}
