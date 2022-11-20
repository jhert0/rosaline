package chess

import (
	"testing"
)

func TestSetBit(t *testing.T) {
	bb := NewBitBoard(0)
	bb.SetBit(1)
	if bb.Value() != 2 {
		t.Fatalf("%v: expected: 2 got %v", t.Name(), bb.Value())
	}
}

func TestClearBit(t *testing.T) {
	bb := NewBitBoard(1)
	bb.ClearBit(0)
	if bb.Value() != 0 {
		t.Fatalf("%v: expected: 0 got %v", t.Name(), bb.Value())
	}
}

func TestBitSet(t *testing.T) {
	bb := NewBitBoard(1)
	if bb.BitSet(1) {
		t.Fatalf("%v: expected: true got %v", t.Name(), bb.Value())
	}
}
