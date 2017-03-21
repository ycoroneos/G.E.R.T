package embedded

type ADC_reading struct {
	Channel uint8
	Value   float32
}

type MCP3008_controller struct {
	spi SPI_periph
}

func MakeMCP3008(spi SPI_periph) *MCP3008_controller {
	//32bit frames in mode 0
	spi.Begin(0, 10, 24, 0)
	return &MCP3008_controller{spi}
}

func (mcp *MCP3008_controller) Read(channel uint8) ADC_reading {
	//stuff gets shifted out in reverse
	channel = channel & 0x7
	command := BitReverse32(uint32(0x3<<3|channel) << 12)
	result := float32(mcp.spi.Exchange(command) & 0x3ff)
	return ADC_reading{channel, (result * 5.0) / 1024.0}
}

//from the internet
func BitReverse32(x uint32) uint32 {
	x = (x&0x55555555)<<1 | (x&0xAAAAAAAA)>>1
	x = (x&0x33333333)<<2 | (x&0xCCCCCCCC)>>2
	x = (x&0x0F0F0F0F)<<4 | (x&0xF0F0F0F0)>>4
	x = (x&0x00FF00FF)<<8 | (x&0xFF00FF00)>>8
	return (x&0x0000FFFF)<<16 | (x&0xFFFF0000)>>16
}
