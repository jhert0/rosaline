package search

import "testing"

const hash uint64 = 16046803855257665054

func TestPush(t *testing.T) {
	table := newDrawTable()

	table.Push(hash)
	retrieved := table.hashes[0]
	if retrieved != hash {
		t.Fatalf("%s: expected hash '%v' got '%v'", t.Name(), hash, retrieved)
	}
}

func TestPop(t *testing.T) {
	table := newDrawTable()
	table.Push(hash)
	popped, ok := table.Pop()
	t.Log("ok:", ok)
	if popped != hash {
		t.Fatalf("%s: expected hash '%v' got '%v'", t.Name(), hash, popped)
	}
}

func TestIsRepeat(t *testing.T) {
	table := newDrawTable()

	table.Push(hash)
	table.Push(hash)

	if table.IsRepeat(hash) != false {
		t.Fatalf("%s: expected not be a draw but one was found", t.Name())
	}

	table.Push(hash)

	if table.IsRepeat(hash) != true {
		t.Fatalf("%s: expected to be a draw but one was not found", t.Name())
	}
}
