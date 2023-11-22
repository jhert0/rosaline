package search

const (
	maxDrawTableSize = 200
)

type drawTable struct {
	hashes []uint64
	index  int
}

// newDrawTable creates a new draw table.
func newDrawTable() drawTable {
	return drawTable{
		hashes: make([]uint64, maxDrawTableSize),
		index:  0,
	}
}

// Push addes a new hash to the draw table.
func (t *drawTable) Push(hash uint64) {
	t.hashes[t.index] = hash
	t.index++
}

// Pop removes the last hash from the draw table.
//
// If the stack is not empty it will return the removed value and true.
// If the stack is empty it will return 0 and false.
func (t *drawTable) Pop() (uint64, bool) {
	if t.index == 0 {
		return 0, false
	}

	value := t.hashes[t.index]
	t.index--

	return value, true
}

func (t *drawTable) Clear() {
	t.index = 0
}

// IsRepeat returns whether the position has been encoutered before
func (t drawTable) IsRepeat(searchHash uint64) bool {
	count := 0
	for i := 0; i < t.index; i++ {
		hash := t.hashes[i]
		if hash == searchHash {
			count++
			if count >= 2 {
				return true
			}
		}
	}

	return false
}
