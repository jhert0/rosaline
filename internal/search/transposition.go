package search

import (
	"fmt"
	"rosaline/internal/chess"
	"unsafe"
)

type NodeType uint8

const (
	ExactNode NodeType = iota
	UpperNode
	LowerNode
)

func (t NodeType) String() string {
	switch t {
	case ExactNode:
		return "Exact"
	case UpperNode:
		return "Upper"
	case LowerNode:
		return "Lower"
	}

	panic(fmt.Sprintf("unknown NodeType '%d' encountered", t))
}

type TableEntry struct {
	Type  NodeType
	Move  chess.Move
	Score int
	Depth int
	Age   int
}

var emptyEntry = TableEntry{}

const (
	entrySize = int(unsafe.Sizeof(emptyEntry))

	kb           = 1024
	mb           = kb * kb
	maxTableSize = (64 * mb) / entrySize
)

// NewTableEntry creates a new TableEntry.
func NewTableEntry(nodeType NodeType, move chess.Move, score, depth int, age int) TableEntry {
	return TableEntry{
		Type:  nodeType,
		Move:  move,
		Score: score,
		Depth: depth,
		Age:   age,
	}
}

func (e TableEntry) String() string {
	return fmt.Sprintf("<Entry: type: %s move: %s score: %d depth: %d>", e.Type, e.Move, e.Score, e.Depth)
}

type TranspositionTable struct {
	table  map[uint64]TableEntry
	hits   int
	misses int
}

// NewTranspositionTable creates a new TranspositionTable.
func NewTranspositionTable() TranspositionTable {
	return TranspositionTable{
		table:  make(map[uint64]TableEntry),
		hits:   0,
		misses: 0,
	}
}

// Insert adds a new entry to the table.
func (t *TranspositionTable) Insert(hash uint64, entry TableEntry) {
	if len(t.table) >= maxTableSize { // TODO: look into replacing old positions instead of clearing table
		clear(t.table)
	}

	t.table[hash] = entry
}

// Remove removes an entry from the table.
func (t *TranspositionTable) Remove(hash uint64) {
	delete(t.table, hash)
}

// Get retreives the entry that corresponds to the given hash.
func (t *TranspositionTable) Get(hash uint64) (TableEntry, bool) {
	value, ok := t.table[hash]

	if ok {
		t.hits++
	} else {
		t.misses++
	}

	return value, ok
}

// Size returns the size of the table.
func (t TranspositionTable) Size() int {
	return len(t.table)
}

// Hits returns the number times a position has been found in the table.
func (t TranspositionTable) Hits() int {
	return t.hits
}

// Misses returns the number times a position wa not found in the table.
func (t TranspositionTable) Misses() int {
	return t.misses
}

func (t TranspositionTable) List() {
	entries := map[NodeType]int{
		ExactNode: 0,
		UpperNode: 0,
		LowerNode: 0,
	}

	for hash, entry := range t.table {
		fmt.Printf("%d: %s\n", hash, entry)
		entries[entry.Type]++
	}

	for key, value := range entries {
		fmt.Printf("%s: %d\n", key, value)
	}
	fmt.Println("# of entries:", len(t.table))
}

// ResetCounters resets the hits and misses counters.
func (t *TranspositionTable) ResetCounters() {
	t.hits = 0
	t.misses = 0
}

// Clear clears the table and resets the hits and misses counters.
func (t *TranspositionTable) Clear() {
	clear(t.table)
	t.ResetCounters()
}
