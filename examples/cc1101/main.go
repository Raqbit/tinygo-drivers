package main

import (
	"encoding/hex"
	"machine"
	"tinygo.org/x/drivers/cc1101"
)

func main() {
	err := machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 433.2e6,
		LSBFirst:  false,
		Mode:      0,
	})

	if err != nil {
		println(err.Error())
	}

	dev := cc1101.New(machine.SPI0, machine.PB2, machine.PB4)

	// Configure device
	if err = dev.Configure(868, 48.0, 48.0, 135.0, 10, 16); err != nil {
		panic(err.Error())
	}

	val, err := dev.GetTestValue()

	if err != nil {
		panic(err.Error())
	}

	Print("Test value (0x59): 0x")
	Println(hex.EncodeToString([]byte{val}))
}

func Print(msg string) {
	machine.UART0.Write([]byte(msg))
}

func Println(msg string) {
	machine.UART0.Write([]byte(msg + "\r\n"))
}
