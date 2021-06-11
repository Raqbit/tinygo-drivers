package cc1101

import (
	"errors"
	"machine"
	"math"
	"time"
	"tinygo.org/x/drivers"
)

const (
	maxPacketLen = 255
	crystalFreq  = 26.0 // Crystal frequency in Mhz
)

// Transfer types
const (
	READ_MODE  = 1 << 7
	BURST_MODE = 1 << 6
)

var (
	errFreq          = errors.New("invalid frequency")
	errPower         = errors.New("invalid output power")
	errRange         = errors.New("invalid range")
	errPacketTooLong = errors.New("packet too long")
	errPreambleLen   = errors.New("invalid preamble length")
	errSyncWord      = errors.New("invalid sync word")
)

type Device struct {
	bus  drivers.SPI
	cs   machine.Pin
	sdo  machine.Pin
	gdo0 machine.Pin

	modulation      Modulation
	freq            float32
	power           int8
	bitRate         float32
	packetLen       uint8
	packetLenConfig PacketLenCfg
}

func New(bus drivers.SPI, cs, miso machine.Pin) Device {
	return Device{
		bus:        bus,
		cs:         cs,
		sdo:        miso,
		modulation: Mod2Fsk,
		//gdo0: gdo0,
	}
}

// Configure configures the CC1101 and all pins used.
func (d Device) Configure(freq, dr, freqDev, rxBw float32, power int8, preambleLen uint8) error {
	d.cs.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Reset cc1101
	d.reset()

	// Set carrier frequency
	if err := d.SetFrequency(freq); err != nil {
		return err
	}

	// Set data rate
	if err := d.SetBitRate(dr); err != nil {
		return err
	}

	// Set receive bandwidth
	if err := d.SetRxBandwidth(rxBw); err != nil {
		return err
	}

	// Set frequency deviation
	if err := d.SetFrequencyDeviation(freqDev); err != nil {
		return err
	}

	// Set default output power
	if err := d.SetOutputPower(power); err != nil {
		return err
	}

	// Set variable packet length mode to maximum
	if err := d.SetVariablePacketLengthMode(maxPacketLen); err != nil {
		return err
	}

	// Set preamble length
	if err := d.SetPreambleLength(preambleLen); err != nil {
		return err
	}

	// Set default modulation
	if err := d.SetModulation(Mod2Fsk); err != nil {
		return err
	}

	// Set data encoding
	if err := d.SetEncoding(false, false); err != nil {
		return err
	}

	// Set sync word
	if err := d.SetSyncWord([2]uint8{0x12, 0xAD}, 0, false); err != nil {
		return err
	}

	// Flush RX FIFO
	if _, err := d.sendCommand(reg_SFRX); err != nil {
		return err
	}

	// Flush TX FIFO
	if _, err := d.sendCommand(reg_SFTX); err != nil {
		return err
	}

	return nil
}

func (d Device) cSelect() {
	d.cSelect()
}

func (d Device) cDeselect() {
	d.cDeselect()
}

// reset executes a manual power-up sequence to get the
// CC1101 into a known state, for instance when the
// power-on-reset did not function properly because the
// power supply did not comply to given specifications.
func (d Device) reset() {
	d.cDeselect()
	time.Sleep(time.Microsecond * 5)
	d.cSelect()
	time.Sleep(time.Microsecond * 10)
	d.cDeselect()
	time.Sleep(time.Microsecond * 41)
	d.cSelect()

	d.waitSDO()
	_, _ = d.bus.Transfer(reg_SRES)
	d.waitSDO()

	d.cDeselect()

	// TODO: set registers
}

// waitSDO waits until the SDO SPI pin goes low,
// which is used by the CC1101 to indicate it is ready
// for an SPI transfer.
func (d Device) waitSDO() {
	for d.sdo.Get() {
	}
}

// sendCommand strobes a C1101 register.
func (d Device) sendCommand(addr uint8) (uint8, error) {
	d.cSelect()
	defer d.cDeselect()

	d.waitSDO()

	res, err := d.bus.Transfer(addr)

	return res, err
}

// readRegister reads from a CC1101 register.
func (d Device) readRegister(addr uint8) (uint8, error) {
	d.cSelect()
	defer d.cDeselect()

	d.waitSDO()

	data := []byte{addr | READ_MODE, 0}

	err := d.bus.Tx(data, data)

	return data[1], err
}

// burstReadRegister reads from multiple CC1101 registers at once.
func (d Device) burstReadRegister(addr uint8, buf []byte) error {
	d.cSelect()
	defer d.cDeselect()

	d.waitSDO()

	// Transfer address
	if _, err := d.bus.Transfer(addr | READ_MODE | BURST_MODE); err != nil {
		return err
	}

	// Read into buffer
	return d.bus.Tx(nil, buf)
}

// writeRegister writes to a CC1101 register.
func (d Device) writeRegister(addr, value uint8) error {
	d.cSelect()
	defer d.cDeselect()

	d.waitSDO()

	data := []byte{addr, value}

	return d.bus.Tx(data, data)
}

// TODO: cache registers instead of read-then-write?
func (d Device) writePartialRegister(addr, value, msb, lsb uint8) error {
	if msb > 7 || lsb > 7 || lsb > msb {
		return errRange
	}

	// Read current register value
	curr, err := d.readRegister(addr)

	if err != nil {
		return err
	}

	// Mask for bits we want to write
	var mask uint8 = ^((0b11111111 << (msb + 1)) | (0b11111111 >> (8 - lsb)))

	// Create new register value from current & new with mask
	newValue := (curr & ^mask) | (value & mask)

	// Write new register value
	return d.writeRegister(addr, newValue)
}

// burstWriteRegister writes to multiple CC1101 registers at once.
func (d Device) burstWriteRegister(addr uint8, buf []byte) error {
	d.cSelect()
	defer d.cDeselect()

	d.waitSDO()

	// Transfer address
	if _, err := d.bus.Transfer(addr | BURST_MODE); err != nil {
		return err
	}

	// Transfer buffer
	return d.bus.Tx(buf, nil)
}

// GetTestValue TODO: remove
func (d Device) GetTestValue() (uint8, error) {
	return d.readRegister(reg_FSTEST)
}

func (d Device) SetFrequency(freq float32) error {
	if !(((freq > 300.0) && (freq < 348.0)) ||
		((freq > 387.0) && (freq < 464.0)) ||
		((freq > 779.0) && (freq < 928.0))) {
		return errFreq
	}

	if _, err := d.sendCommand(reg_SIDLE); err != nil {
		return err
	}

	var frf = uint32((freq * (1 << 16)) / crystalFreq)
	if err := d.writeRegister(reg_FREQ2, uint8((frf&0xFF0000)>>16)); err != nil {
		return err
	}
	if err := d.writeRegister(reg_FREQ1, uint8((frf&0x00FF00)>>8)); err != nil {
		return err
	}
	if err := d.writeRegister(reg_FREQ0, uint8(frf&0x0000FF)); err != nil {
		return err
	}

	// Store new frequency
	d.freq = freq

	// Update TX power for new frequency
	if err := d.SetOutputPower(d.power); err != nil {
		return err
	}

	return nil
}

func (d Device) SetOutputPower(power int8) error {
	var pwrIdx uint8

	switch power {
	case -30:
		pwrIdx = 0
	case -20:
		pwrIdx = 1
	case -15:
		pwrIdx = 2
	case -10:
		pwrIdx = 3
	case 0:
		pwrIdx = 4
	case 5:
		pwrIdx = 5
	case 7:
		pwrIdx = 6
	case 10:
		pwrIdx = 7
	default:
		return errPower
	}

	var freqIdx uint8

	if d.freq < 374.0 {
		freqIdx = 0
	} else if d.freq < 650.5 {
		freqIdx = 1
	} else if d.freq < 891.5 {
		freqIdx = 2
	} else {
		freqIdx = 3
	}

	// Table with values for PATABLE registers
	// the configured frequency is used to get the column,
	// the given power is used to get the the row
	//
	// Values from Table 39 in datasheet
	paTable := [][]uint8{
		{0x12, 0x12, 0x03, 0x03},
		{0x0D, 0x0E, 0x0F, 0x0E},
		{0x1C, 0x1D, 0x1E, 0x1E},
		{0x34, 0x34, 0x27, 0x27},
		{0x51, 0x60, 0x50, 0x8E},
		{0x85, 0x84, 0x81, 0xCD},
		{0xCB, 0xC8, 0xCB, 0xC7},
		{0xC2, 0xC0, 0xC2, 0xC0},
	}

	// Get power value from PA table
	rawPower := paTable[pwrIdx][freqIdx]

	var err error

	if d.modulation == ModAskOok {
		// Amplitude modulation, off (0x00) or on (full power)
		paValues := [2]uint8{0x00, rawPower}
		err = d.burstWriteRegister(reg_PATABLE, paValues[:])
	} else {
		// Frequency modulation, use full power when transmitting
		err = d.writeRegister(reg_PATABLE, rawPower)
	}

	// Save new power value
	d.power = power

	return err
}

func (d Device) SetBitRate(bitRate float32) error {
	if !(bitRate >= 0.025 && bitRate <= 600) {
		return errRange
	}

	if _, err := d.sendCommand(reg_SIDLE); err != nil {
		return err
	}

	exp, mant := getExpMant(bitRate*1000, 256, 28, 14)

	if err := d.writePartialRegister(reg_MDMCFG4, exp, 3, 0); err != nil {
		return err
	}

	if err := d.writeRegister(reg_MDMCFG3, mant); err != nil {
		return err
	}

	d.bitRate = bitRate

	return nil
}

func (d Device) SetRxBandwidth(bw float32) error {
	if !(bw >= 58 && bw <= 812) {
		return errRange
	}

	if _, err := d.sendCommand(reg_SIDLE); err != nil {
		return err
	}

	for e := 3; e >= 0; e-- {
		for m := 3; m >= 0; m-- {
			point := (crystalFreq * 1000000.0) / (8 * (m + 4) * (1 << e))
			if (math.Abs(float64(bw)*1000.0) - point) <= 1000 {
				// set Rx channel filter bandwidth
				var value = uint8((e << 6) | (m << 4))

				return d.writePartialRegister(reg_MDMCFG4, value, 7, 4)
			}
		}
	}

	panic("Could not find exponent for rx bandwidth")
}

func (d Device) SetFrequencyDeviation(freqDev float32) error {
	if freqDev < 0.0 {
		freqDev = 1.587
	}

	if !(freqDev >= 1.587 && freqDev <= 380.8) {
		return errRange
	}

	if _, err := d.sendCommand(reg_SIDLE); err != nil {
		return err
	}

	exp, mant := getExpMant(freqDev*1000, 8, 17, 7)

	if err := d.writePartialRegister(reg_DEVIATN, exp<<4, 6, 4); err != nil {
		return err
	}

	return d.writePartialRegister(reg_DEVIATN, mant, 2, 0)
}

func (d Device) setPacketLength(mode PacketLenCfg, len uint8) error {
	if len > maxPacketLen {
		return errPacketTooLong
	}

	if err := d.writePartialRegister(reg_PKTCTRL0, uint8(mode), 1, 0); err != nil {
		return err
	}

	if err := d.writeRegister(reg_PKTLEN, len); err != nil {
		return err
	}

	d.packetLen = len
	d.packetLenConfig = mode

	return nil
}

func (d Device) SetFixedPacketLengthMode(len uint8) error {
	return d.setPacketLength(PktLenFixed, len)
}

func (d Device) SetVariablePacketLengthMode(maxLen uint8) error {
	return d.setPacketLength(PktlenVariable, maxLen)
}

func (d Device) SetPreambleLength(preambleLen uint8) error {
	var value PreambleLen

	switch preambleLen {
	case 16:
		value = Preamble2
	case 24:
		value = Preamble3
	case 32:
		value = Preamble4
	case 48:
		value = Preamble6
	case 64:
		value = Preamble8
	case 96:
		value = Preamble12
	case 128:
		value = Preamble16
	case 192:
		value = Preamble24
	default:
		return errPreambleLen
	}

	return d.writePartialRegister(reg_MDMCFG1, uint8(value), 6, 4)
}

func (d Device) SetModulation(mod Modulation) error {
	if _, err := d.sendCommand(reg_SIDLE); err != nil {
		return err
	}

	if mod == ModAskOok {
		// We are switching to amplitude modulation, set PA power setting
		if err := d.writePartialRegister(reg_FREND0, 1, 2, 0); err != nil {
			return err
		}
	} else if d.modulation == ModAskOok {
		// We are switching away from amplitude modulation, reset PA power setting
		if err := d.writePartialRegister(reg_FREND0, 0, 2, 0); err != nil {
			return err
		}
	}

	// Store new modulation
	d.modulation = mod

	// Update modulation
	if err := d.writePartialRegister(reg_MDMCFG2, uint8(mod), 6, 4); err != nil {
		return err
	}

	// Update output power for new modulation
	return d.SetOutputPower(d.power)
}

func (d Device) SetEncoding(manchesterEn bool, whiteningEn bool) error {
	if _, err := d.sendCommand(reg_SIDLE); err != nil {
		return err
	}

	manchester := ManchesterOff
	if manchesterEn {
		manchester = ManchesterOn
	}

	whitening := WhiteDataOff
	if whiteningEn {
		whitening = WhiteDataOn
	}

	// Configure manchester coding
	if err := d.writePartialRegister(reg_MDMCFG2, uint8(manchester), 3, 3); err != nil {
		return err
	}

	// Configure data whitening
	return d.writePartialRegister(reg_PKTCTRL0, uint8(whitening), 6, 6)
}

func (d Device) SetSyncWord(syncWord [2]uint8, maxErrBits uint8, requireCarrierSense bool) error {
	if maxErrBits > 1 {
		return errSyncWord
	}

	if syncWord[0] == 0x00 || syncWord[1] == 0x00 {
		return errSyncWord
	}

	if err := d.EnableSyncWordFiltering(maxErrBits, requireCarrierSense); err != nil {
		return err
	}

	if err := d.writeRegister(reg_SYNC1, syncWord[0]); err != nil {
		return err
	}

	return d.writeRegister(reg_SYNC0, syncWord[1])
}

func (d Device) EnableSyncWordFiltering(maxErrBits uint8, requireCarrierSense bool) error {
	switch maxErrBits {
	case 0:
		val := SyncMode1616

		if requireCarrierSense {
			val = SyncMode1616Thr
		}

		return d.writePartialRegister(reg_MDMCFG2, uint8(val), 2, 0)
	case 1:
		val := SyncMode1516

		if requireCarrierSense {
			val = SyncMode1516Thr
		}

		return d.writePartialRegister(reg_MDMCFG2, uint8(val), 2, 0)
	default:
		return errSyncWord
	}
}

// An approximation of the formula on page 35 of CC1101 data sheet
// Taken from https://github.com/jgromes/RadioLib/blob/251dd438a028167895b3cda35f107da90d9fc1c7/src/modules/CC1101/CC1101.cpp#L903
func getExpMant(target float32, mantOffset uint16, divExp uint8, expMax int8) (exp uint8, mant uint8) {
	origin := (float32(mantOffset) * float32(crystalFreq) * 1000000.0) / (1 << divExp)
	for e := expMax; e >= 0; e-- {
		intervalStart := (1 << e) * origin

		if target >= intervalStart {
			exp = uint8(e)

			stepSize := intervalStart / mantOffset

			mant = (target - intervalStart) / stepSize
			return
		}
	}

	return
}
