package chess

type Color uint8

const (
	NoColor Color = 0
	White   Color = 1
	Black   Color = 2
)

// String returns the string representation of the Color.
func (c Color) String() string {
	if c == White {
		return "White"
	}

	if c == Black {
		return "Black"
	}

	return "No Color"
}
