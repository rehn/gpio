package gpio

import "time"

type Button struct {
	gpio           Gpio
	RepeatInterval int64
	KeyPress       func()
	KeyUp          func()
}

func NewButton(gpioPin int, repeatIterval int64) Button {
	g := NewGpio("in", gpioPin)
	btn := Button{gpio: g, RepeatInterval: repeatIterval, KeyPress: func() {}(), KeyUp: func() {}()}
	go btn.buttonWatcher()
	return btn
}

func (b *Button) buttonWatcher() {
	currentValue := b.gpio.GetValue()
	for {
		newValue := b.gpio.GetValue()
		if currentValue != newValue || newValue == true {
			if newValue == true {
				b.KeyPress()
				var d time.Duration = time.Duration(time.Duration(b.RepeatInterval-10) * time.Millisecond)
				time.Sleep(d)
			} else {
				b.KeyUp()
			}
			currentValue = newValue
		}
		time.Sleep(10 * time.Millisecond)
	}
}
