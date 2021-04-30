package cc1101

// PATRABLE & FIFO's
const (
	reg_PATABLE = 0x3E // Power amplifier table
	reg_FIFO    = 0x3F // FIFO access
)

// Configuration registers
const (
	reg_IOCFG2   = 0x00 // GDO2 Output Pin Configuration
	reg_IOCFG1   = 0x01 // GDO1 Output Pin Configuration
	reg_IOCFG0   = 0x02 // GDO0 Output Pin Configuration
	reg_FIFOTHR  = 0x03 // RX FIFO and TX FIFO Thresholds
	reg_SYNC1    = 0x04 // Sync Word, High Byte
	reg_SYNC0    = 0x05 // Sync Word, Low Byte
	reg_PKTLEN   = 0x06 // Packet Length
	reg_PKTCTRL1 = 0x07 // Packet Automation Control
	reg_PKTCTRL0 = 0x08 // Packet Automation Control
	reg_ADDR     = 0x09 // Device Address
	reg_CHANNR   = 0x0A // Channel Number
	reg_FSCTRL1  = 0x0B // Frequency Synthesizer Control
	reg_FSCTRL0  = 0x0C // Frequency Synthesizer Control
	reg_FREQ2    = 0x0D // Frequency Control Word, High Byte
	reg_FREQ1    = 0x0E // Frequency Control Word, Middle Byte
	reg_FREQ0    = 0x0F // Frequency Control Word, Low Byte
	reg_MDMCFG4  = 0x10 // Modem Configuration
	reg_MDMCFG3  = 0x11 // Modem Configuration
	reg_MDMCFG2  = 0x12 // Modem Configuration
	reg_MDMCFG1  = 0x13 // Modem Configuration
	reg_MDMCFG0  = 0x14 // Modem Configuration
	reg_DEVIATN  = 0x15 // Modem Deviation Setting
	reg_MCSM2    = 0x16 // Main Radio Control State Machine Configuration
	reg_MCSM1    = 0x17 // Main Radio Control State Machine Configuration
	reg_MCSM0    = 0x18 // Main Radio Control State Machine Configuration
	reg_FOCCFG   = 0x19 // Frequency Offset Compensation Configuration
	reg_BSCFG    = 0x1A // Bit Synchronization Configuration
	reg_AGCCTRL2 = 0x1B // AGC Control
	reg_AGCCTRL1 = 0x1C // AGC Control
	reg_AGCCTRL0 = 0x1D // AGC Control
	reg_WOREVT1  = 0x1E // High Byte Event0 Timeout
	reg_WOREVT0  = 0x1F // Low Byte Event0 Timeout
	reg_WORCTRL  = 0x20 // Wake On Radio Control
	reg_FREND1   = 0x21 // Front End RX Configuration
	reg_FREND0   = 0x22 // Front End TX Configuration
	reg_FSCAL3   = 0x23 // Frequency Synthesizer Calibration
	reg_FSCAL2   = 0x24 // Frequency Synthesizer Calibration
	reg_FSCAL1   = 0x25 // Frequency Synthesizer Calibration
	reg_FSCAL0   = 0x26 // Frequency Synthesizer Calibration
	reg_RCCTRL1  = 0x27 // RC Oscillator Configuration
	reg_RCCTRL0  = 0x28 // RC Oscillator Configuration

	// These are not preserved in SLEEP state

	reg_FSTEST  = 0x29 // Frequency Synthesizer Calibration Control
	reg_PTEST   = 0x2A // Production Test
	reg_AGCTEST = 0x2B // AGC Test
	reg_TEST2   = 0x2C // Various Test Settings
	reg_TEST1   = 0x2D // Various Test Settings
	reg_TEST0   = 0x2E // Various Test Settings
)

// Command Strobe Registers
const (
	reg_SRES = 0x30 // Reset chip

	// Enable and calibrate frequency synthesizer (if MCSM0.FS_AUTOCAL=1).
	// If in RX (with CCA): Go to a wait state where only the synthesizer
	// is running (for quick RX / TX turnaround).
	reg_SFSTXON = 0x31

	reg_SXOFF = 0x32 // Turn off crystal oscillator.

	// Calibrate frequency synthesizer and turn it off.
	// SCAL can be strobed from IDLE mode without setting manual calibration mode.
	reg_SCAL = 0x33

	reg_SRX = 0x34 // Enable RX. Perform calibration first if coming from IDLE and MCSM0.FS_AUTOCAL=1.

	// In IDLE state: Enable TX. Perform calibration first
	// if MCSM0.FS_AUTOCAL=1.
	// If in RX state and CCA is enabled: Only go to TX if channel is clear.
	reg_STX = 0x35

	reg_SIDLE = 0x36 // Exit RX / TX, turn off frequency synthesizer and exit Wake-On-Radio mode if applicable.

	// Start automatic RX polling sequence (Wake-on-Radio) as described in Section 19.5 if WORCTRL.RC_PD=0.
	reg_SWOR = 0x38

	reg_SPWD    = 0x39 // Enter power down mode when CSn goes high.
	reg_SFRX    = 0x3A // Flush the RX FIFO buffer. Only issue SFRX in IDLE or RXFIFO_OVERFLOW states.
	reg_SFTX    = 0x3B // Flush the TX FIFO buffer. Only issue SFTX in IDLE or TXFIFO_UNDERFLOW states.
	reg_SWORRST = 0x3C // Reset real time clock to Event1 value.
	reg_SNOP    = 0x3D // No operation. May be used to get access to the chip status byte.
)

// Status Registers
const (
	reg_PARTNUM        = 0xF0 // Chip ID
	reg_VERSION        = 0xF1 // Chip ID
	reg_FREQEST        = 0xF2 // Frequency Offset Estimate from Demodulator
	reg_LQI            = 0xF3 // Demodulator Estimate for Link Quality
	reg_RSSI           = 0xF4 // Received Signal Strength Indication
	reg_MARCSTATE      = 0xF5 // Main Radio Control State Machine State
	reg_WORTIME1       = 0xF6 // High Byte of WOR Time
	reg_WORTIME0       = 0xF7 // Low Byte of WOR Time
	reg_PKTSTATUS      = 0xF8 // Current GDOx Status and Packet Status
	reg_VCO_VC_DAC     = 0xF9 // Current Setting from PLL Calibration Module
	reg_TXBYTES        = 0xFA // Underflow and Number of Bytes
	reg_RXBYTES        = 0xFB // Overflow and Number of Bytes
	reg_RCCTRL1_STATUS = 0xFC // Last RC Oscillator Calibration Result
	reg_RCCTRL0_STATUS = 0xFD // Last RC Oscillator Calibration Result
)
