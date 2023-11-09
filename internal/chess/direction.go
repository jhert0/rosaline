package chess

import "fmt"

type direction int8

// These directions are from white's perspective.
const (
	north direction = 8
	south direction = -8

	east direction = 1
	west direction = -1

	northwest direction = north + west
	northeast direction = north + east

	southwest direction = south + west
	southeast direction = south + east
)

var directions = []direction{
	north,
	south,
	east,
	west,

	northeast,
	northwest,
	southeast,
	southwest,
}

func (d direction) String() string {
	switch d {
	case north:
		return "North"
	case south:
		return "South"
	case east:
		return "East"
	case west:
		return "West"
	case north + east:
		return "Northeast"
	case north + west:
		return "Northwest"
	case south + east:
		return "Southeast"
	case south + west:
		return "Southwest"
	}

	panic(fmt.Sprintf("Unknown direction: %d", d))
}

func (d direction) rayIndex() int {
	switch d {
	case north:
		return 0
	case south:
		return 1
	case east:
		return 2
	case west:
		return 3
	case northwest:
		return 4
	case northeast:
		return 5
	case southwest:
		return 6
	case southeast:
		return 7
	}

	panic(fmt.Sprintf("rayIndex: unknown direction encountered: %d", d))
}
