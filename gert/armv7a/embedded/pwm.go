package embedded

const (
	clk_ipg          = 1
	clk_ipg_highfreq = 2
	clk_ipg_32k      = 3
)

type khz uint32

type PWM_regs struct {
	CR  uint32
	SR  uint32
	IR  uint32
	SAR uint32
	PR  uint32
	CNR uint32
}

type PWM_pin struct {
	name   string
	alt    uint8
	muxctl *uint32
	padctl *uint32
}

type PWM_periph struct {
	regs           *PWM_regs
	output         PWM_pin
	mode           uint8
	frequency      uint8
	channel_select uint8
}

func (pwm *PWM_periph) Begin(freq khz) {
	//section 51.5 of the DQRM

	//disable pwm
	pwm.regs.CR = 0

	//set up CR
	//pwm goes high on a rollover
	//pwm clk source is ipg_clk (system clock)
	pwm.regs.CR |= clk_ipg << 16

	//set the prescalar, which effectively sets the switching frequency
	prescale := freq & 0xFFF
	pwm.regs.CR |= uint32(prescale) << 4

	//turn of all pwm interrupts
	pwm.regs.IR = 0

	//clear all the things in the status register
	pwm.regs.SR = 0xFF

	//set the period to the max, we will use the prescalar to change the period
	pwm.regs.PR = 0xFFFF

	//enable pwm
	pwm.regs.CR |= 1
}

func (pwm *PWM_periph) Stop() {
	pwm.regs.CR &= ^uint32(1)
}

func (pwm *PWM_periph) SetFreq(freq khz) {
	//we can technically change the divider while its running
	cr := pwm.regs.CR
	cr &= ^uint32(0xFFF)
	cr |= (uint32(freq) & 0xFFF)
	pwm.regs.CR = cr
}

func (pwm *PWM_periph) SetDuty(dutycycle float32) {
	realduty := uint32(0xFFFF * dutycycle)
	pwm.regs.SAR = realduty
}
