package embedded

type MDD10A_controller struct {
	PWM1 PWM_periph
	DIR1 GPIO_pin

	PWM2 PWM_periph
	DIR2 GPIO_pin
}

func MakeMDD10A(pwm1, pwm2 PWM_periph, dir1, dir2 GPIO_pin) *MDD10A_controller {
	pwm1.Begin(0x100)
	pwm1.SetDuty(0.0)
	pwm2.Begin(0x100)
	pwm2.SetDuty(0.0)
	dir1.SetOutput()
	dir2.SetOutput()
	return &MDD10A_controller{pwm1, dir1, pwm2, dir2}
}

func (c *MDD10A_controller) move(speed1, speed2 float32, dir1, dir2 bool) {
	if dir1 {
		c.DIR1.SetHI()
	} else {
		c.DIR1.SetLO()
	}

	if dir2 {
		c.DIR2.SetHI()
	} else {
		c.DIR2.SetLO()
	}

	c.PWM1.SetDuty(speed1)
	c.PWM2.SetDuty(speed2)
}

func (c *MDD10A_controller) Forward(speed float32) {
	c.move(speed, speed, true, true)
}

func (c *MDD10A_controller) Backward(speed float32) {
	c.move(speed, speed, false, false)
}

func (c *MDD10A_controller) TurnLeft(speed float32) {
	c.move(speed, speed, false, true)
}

func (c *MDD10A_controller) TurnRight(speed float32) {
	c.move(speed, speed, true, false)
}

func (c *MDD10A_controller) Stop() {
	c.move(0.0, 0.0, true, true)
}
