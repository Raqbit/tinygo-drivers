package cc1101

import (
	"machine"
	"time"
	"tinygo.org/x/drivers"
)

type CarrierFreq uint8

// Transfer types
const (
	WRITE_SINGLE_BYTE = 0x00
	WRITE_BURST       = 0x40
	READ_SINGLE_BYTE  = 0x80
	READ_BURST        = 0xC0
)

type Device struct {
	bus  drivers.SPI
	cs   machine.Pin
	sdo  machine.Pin
	gdo0 machine.Pin
}

func New(bus drivers.SPI, cs, miso machine.Pin) Device {
	return Device{
		bus: bus,
		cs:  cs,
		sdo: miso,
		//gdo0: gdo0,
	}
}

// Configure configures the CC1101 and all pins used.
func (d Device) Configure() {
	d.cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.sdo.Configure(machine.PinConfig{Mode: machine.PinInput})
	d.cs.High()

	d.reset()
}

// reset executes a manual power-up sequence to get the
// CC1101 into a known state, for instance when the
// power-on-reset did not function properly because the
// power supply did not comply to given specifications.
func (d Device) reset() {
	d.cs.High()
	time.Sleep(time.Microsecond * 5)
	d.cs.Low()
	time.Sleep(time.Microsecond * 10)
	d.cs.High()
	time.Sleep(time.Microsecond * 41)
	d.cs.Low()

	d.waitSDO()
	_, _ = d.bus.Transfer(reg_SRES)
	d.waitSDO()

	d.cs.High()
}

// waitSDO waits until the SDO SPI pin goes low,
// which is used by the CC1101 to indicate it is ready
// for an SPI transfer.
func (d Device) waitSDO() {
	for d.sdo.Get() {
	}
}

// strobeRegister strobes a C1101 register.
func (d Device) strobeRegister(addr uint8) uint8 {
	d.cs.Low()
	defer d.cs.High()

	d.waitSDO()

	res, _ := d.bus.Transfer(addr)

	return res
}

// readRegister reads from a CC1101 register.
func (d Device) readRegister(addr uint8) uint8 {
	d.cs.Low()
	defer d.cs.High()

	d.waitSDO()

	data := []byte{addr | READ_SINGLE_BYTE, 0}

	_ = d.bus.Tx(data, data)

	return data[1]
}

// burstReadRegister reads from multiple CC1101 registers at once.
func (d Device) burstReadRegister(addr uint8, buf []byte) {
	d.cs.Low()
	defer d.cs.High()

	d.waitSDO()

	// Transfer address
	_, _ = d.bus.Transfer(addr | READ_BURST)

	// Read into buffer
	_ = d.bus.Tx(nil, buf)
}

// writeRegister writes to a CC1101 register.
func (d Device) writeRegister(addr uint8, value uint8) {
	d.cs.Low()
	defer d.cs.High()

	d.waitSDO()

	data := []byte{addr | WRITE_SINGLE_BYTE, value}

	_ = d.bus.Tx(data, data)
}

// burstWriteRegister writes to multiple CC1101 registers at once.
func (d Device) burstWriteRegister(addr uint8, buf []byte) {
	d.cs.Low()
	defer d.cs.High()

	d.waitSDO()

	// Transfer address
	_, _ = d.bus.Transfer(addr | WRITE_BURST)

	// Transfer buffer
	_ = d.bus.Tx(buf, nil)
}
