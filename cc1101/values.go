package cc1101

type Modulation uint8

// Modulation format
const (
	Mod2Fsk   = Modulation(0b00000000) // 2-FSK (default)
	ModGfsk   = Modulation(0b00010000) // GFSK
	ModAskOok = Modulation(0b00110000) // ASK/OOK
	Mod4Fsk   = Modulation(0b01000000) // 4-FSK
	ModMfsk   = Modulation(0b01110000) // MFSK - only for data rates above 26 kBaud
)

type PacketLenCfg uint8

const (
	PktLenFixed    = PacketLenCfg(0b00000000)
	PktlenVariable = PacketLenCfg(0b00000001)
	PktLenInfinite = PacketLenCfg(0b00000010)
)

type PreambleLen uint8

const (
	Preamble2  = PreambleLen(0b00000000)
	Preamble3  = PreambleLen(0b00010000)
	Preamble4  = PreambleLen(0b00100000)
	Preamble6  = PreambleLen(0b00110000)
	Preamble8  = PreambleLen(0b01000000)
	Preamble12 = PreambleLen(0b01010000)
	Preamble16 = PreambleLen(0b01100000)
	Preamble24 = PreambleLen(0b01110000)
)

type ManchesterEncoding uint8

const (
	ManchesterOff = ManchesterEncoding(0b00000000)
	ManchesterOn  = ManchesterEncoding(0b00001000)
)

type DataWhitening uint8

const (
	WhiteDataOff = DataWhitening(0b00000000)
	WhiteDataOn  = DataWhitening(0b01000000)
)

type SyncMode uint8

const (
	SyncMode1516    = SyncMode(0b00000001)
	SyncMode1616    = SyncMode(0b00000010)
	SyncMode1516Thr = SyncMode(0b00000101)
	SyncMode1616Thr = SyncMode(0b00000110)
)
