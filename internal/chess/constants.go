package chess

import "fmt"

const (
	StartingFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

const (
	FileABB = BitBoard(72340172838076673)
	FileBBB = BitBoard(144680345676153346)
	FileCBB = BitBoard(289360691352306692)
	FileDBB = BitBoard(578721382704613384)
	FileEBB = BitBoard(1157442765409226768)
	FileFBB = BitBoard(2314885530818453536)
	FileGBB = BitBoard(4629771061636907072)
	FileHBB = BitBoard(9259542123273814144)
)

type direction int8

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
