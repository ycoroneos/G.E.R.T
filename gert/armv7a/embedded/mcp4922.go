package embedded

type MCP4922_controller struct {
	spi SPI_periph
}

func MakeMCP4922(spi SPI_periph) MCP4922_controller {
	//16bit frames in mode 0
	spi.Begin(0, 10, 16, 0)
	return MCP4922_controller{spi}
}

func (m MCP4922_controller) Write(data uint16, channel uint8) {
	channel &= 0x1
	data &= 0xFFF
	out := (uint32(channel) << 15) | (uint32(1) << 13) | (uint32(1) << 12) | uint32(data)
	m.spi.Send(out)
}
