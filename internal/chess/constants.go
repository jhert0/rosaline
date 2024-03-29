package chess

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

var FileBitBoards = []BitBoard{
	FileABB,
	FileBBB,
	FileCBB,
	FileDBB,
	FileEBB,
	FileFBB,
	FileGBB,
	FileHBB,
}

const (
	Rank1BB = BitBoard(255)
	Rank2BB = BitBoard(65280)
	Rank3BB = BitBoard(16711680)
	Rank4BB = BitBoard(4278190080)
	Rank5BB = BitBoard(1095216660480)
	Rank6BB = BitBoard(280375465082880)
	Rank7BB = BitBoard(71776119061217280)
	Rank8BB = BitBoard(18374686479671623680)
)
