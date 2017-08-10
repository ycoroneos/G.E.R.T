// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	output PWM_pin
	regs   *PWM_regs
	period khz
	duty   float32
}

func (pwm *PWM_periph) Begin(freq khz) {
	//section 51.5 of the DQRM
	*pwm.output.muxctl = makeGPIOmuxconfig(pwm.output.alt)
	*pwm.output.padctl = makeGPIOpadconfig(1, PULLDOWN_100K, 1, 1, 0, SPEED_FAST, DRIVE_260R, SLEW_FAST)

	//disable pwm
	pwm.regs.CR = 0

	//set up CR
	//pwm goes high on a rollover
	//pwm clk source is ipg_clk (system clock)
	pwm.regs.CR |= clk_ipg << 16

	//set the prescalar, which effectively sets the switching frequency
	//prescale := freq & 0xFFF
	prescale := 0
	pwm.regs.CR |= uint32(prescale) << 4

	//turn of all pwm interrupts
	pwm.regs.IR = 0

	//clear all the things in the status register
	pwm.regs.SR = 0xFF

	//pwm.regs.PR = 0xFFFF
	pwm.regs.PR = uint32(freq)
	pwm.period = freq
	pwm.duty = 0

	//enable pwm
	pwm.regs.CR |= 1
}

func (pwm *PWM_periph) Stop() {
	pwm.regs.CR &= ^uint32(1)
}

func (pwm *PWM_periph) SetFreq(freq khz) {
	//we can technically change the divider while its running
	//cr := pwm.regs.CR
	//cr &= ^uint32(0xFFF)
	//cr |= (uint32(freq) & 0xFFF)
	//pwm.regs.CR = cr
	pwm.period = freq
	pwm.regs.PR = uint32(freq)
	pwm.SetDuty(pwm.duty)
}

func (pwm *PWM_periph) SetDuty(dutycycle float32) {
	pwm.duty = dutycycle
	pwm.regs.SAR = uint32(float32(pwm.period) * pwm.duty)
	//	realduty := uint32(0xFFFF * dutycycle)
	//	pwm.regs.SAR = realduty
}
