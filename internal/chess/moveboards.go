package chess

var knightMoves = [64]BitBoard{
	132096, 329728, 659712, 1319424, 2638848, 5277696, 10489856, 4202496,
	33816580, 84410376, 168886289, 337772578, 675545156, 1351090312, 2685403152, 1075839008,
	8657044482, 21609056261, 43234889994, 86469779988, 172939559976, 345879119952, 687463207072, 275423174720,
	2216203387392, 5531918402816, 11068131838464, 22136263676928, 44272527353856, 88545054707712, 175990581010432, 70506185244672,
	567348067172352, 1416171111120896, 2833441750646784, 5666883501293568, 11333767002587136, 22667534005174272, 45053588738670592, 18049583422636032,
	145241105196122112, 362539804446949376, 725361088165576704, 1450722176331153408, 2901444352662306816, 5802888705324613632, 11533718717099671552, 4620693356194824192,
	288234782788157440, 576469569871282176, 1224997833292120064, 2449995666584240128, 4899991333168480256, 9799982666336960512, 9799982666336960512, 9799982666336960512,
	1128098930098176, 2257297371824128, 4796069720358912, 4796069720358912, 19184278881435648, 38368557762871296, 4679521487814656, 9077567998918656,
}

var kingMoves = [64]BitBoard{
	770, 1797, 3594, 7188, 14376, 28752, 57504, 49216,
	197123, 460039, 920078, 1840156, 3680312, 7360624, 14721248, 12599488,
	50463488, 117769984, 235539968, 471079936, 942159872, 1884319744, 3768639488, 3225468928,
	12918652928, 30149115904, 60298231808, 120596463616, 241192927232, 482385854464, 964771708928, 825720045568,
	3307175149568, 7718173671424, 15436347342848, 30872694685696, 61745389371392, 123490778742784, 246981557485568, 211384331665408,
	846636838289408, 1975852459884544, 3951704919769088, 7903409839538176, 15806819679076352, 31613639358152704, 63227278716305408, 54114388906344448,
	216739030602088448, 505818229730443264, 1011636459460886528, 2023272918921773056, 4046545837843546112, 8093091675687092224, 16186183351374184448, 13853283560024178688,
	144959613005987840, 362258295026614272, 724516590053228544, 1449033180106457088, 2898066360212914176, 5796132720425828352, 11592265440851656704, 4665729213955833856,
}

var rayAttacks [8][65]BitBoard

func init() {
	for square := A1; square <= H8; square++ {
		northBB := BitBoard(0)
		northBB.SetBit(uint64(square))
		northBB.FillNorth()
		northBB.ClearBit(uint64(square))

		rayAttacks[north.rayIndex()][square] = northBB

		southBB := BitBoard(0)
		southBB.SetBit(uint64(square))
		southBB.FillSouth()
		southBB.ClearBit(uint64(square))

		rayAttacks[south.rayIndex()][square] = southBB

		eastBB := BitBoard(0)
		eastBB.SetBit(uint64(square))

		westBB := BitBoard(0)
		westBB.SetBit(uint64(square))

		northEastBB := BitBoard(0)
		northEastBB.SetBit(uint64(square))

		northWestBB := BitBoard(0)
		northWestBB.SetBit(uint64(square))

		southEastBB := BitBoard(0)
		southEastBB.SetBit(uint64(square))

		southWestBB := BitBoard(0)
		southWestBB.SetBit(uint64(square))

		for file := square.File(); file <= 8; file++ {
			eastBB.EastOne()
			northEastBB.NortheastOne()
			southEastBB.SoutheastOne()
		}

		for file := square.File(); file >= 1; file-- {
			westBB.WestOne()
			northWestBB.NorthwestOne()
			southWestBB.SouthwestOne()
		}

		eastBB.ClearBit(uint64(square))
		westBB.ClearBit(uint64(square))
		northEastBB.ClearBit(uint64(square))
		northWestBB.ClearBit(uint64(square))
		southEastBB.ClearBit(uint64(square))
		southWestBB.ClearBit(uint64(square))

		rayAttacks[east.rayIndex()][square] = eastBB
		rayAttacks[west.rayIndex()][square] = westBB
		rayAttacks[northeast.rayIndex()][square] = northEastBB
		rayAttacks[northwest.rayIndex()][square] = northWestBB
		rayAttacks[southeast.rayIndex()][square] = southEastBB
		rayAttacks[southwest.rayIndex()][square] = southWestBB
	}

	for _, direction := range directions {
		rayAttacks[direction.rayIndex()][64] = BitBoard(0)
	}
}

func getPositiveRayAttacks(occupied BitBoard, dir direction, square Square) BitBoard {
	rayIndex := dir.rayIndex()

	attacks := rayAttacks[rayIndex][square]
	blocker := attacks & occupied
	firstBlockingSquare := blocker.Lsb()
	attacks ^= rayAttacks[rayIndex][firstBlockingSquare]
	return attacks
}

func getNegativeRayAttacks(occupied BitBoard, dir direction, square Square) BitBoard {
	rayIndex := dir.rayIndex()

	attacks := rayAttacks[rayIndex][square]
	blocker := attacks & occupied
	firstBlockingSquare := 63 - blocker.Msb()
	if firstBlockingSquare < 0 {
		firstBlockingSquare = 64
	}
	attacks ^= rayAttacks[rayIndex][firstBlockingSquare]
	return attacks
}

func getFileRayAttacks(occupied BitBoard, square Square) BitBoard {
	return getPositiveRayAttacks(occupied, north, square) | getNegativeRayAttacks(occupied, south, square)
}

func getRankRayAttacks(occupied BitBoard, square Square) BitBoard {
	return getPositiveRayAttacks(occupied, east, square) | getNegativeRayAttacks(occupied, west, square)
}

func getDiagonalAttacks(occupied BitBoard, square Square) BitBoard {
	return getPositiveRayAttacks(occupied, northeast, square) | getNegativeRayAttacks(occupied, southwest, square)
}

func getAntiDiagonalAttacks(occupied BitBoard, square Square) BitBoard {
	return getPositiveRayAttacks(occupied, northwest, square) | getNegativeRayAttacks(occupied, southeast, square)
}

func getRookAttacks(occupied BitBoard, square Square) BitBoard {
	return getFileRayAttacks(occupied, square) | getRankRayAttacks(occupied, square)
}

func getBishopAttacks(occupied BitBoard, square Square) BitBoard {
	return getDiagonalAttacks(occupied, square) | getAntiDiagonalAttacks(occupied, square)
}
