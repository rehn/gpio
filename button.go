package gpio

import "time"

type Button struct {
	gpio           Gpio
	buttonDown     chan bool
	buttonUp       chan bool
	repeatInterval int64
}

func NewButton(gpioPin int, repeatIterval int64, repeatAcc int64, repeatStartAcc int64) Button {
	g := NewGpio("in", gpioPin)
	btn := Button{gpio: g, repeatInterval: repeatIterval, buttonDown: make(chan bool), buttonUp: make(chan bool)}
	go btn.buttonWatcher()
	return btn
}

func (b *Button) buttonWatcher() {
	currentValue := b.gpio.getValue()
	for {
		newValue := b.gpio.getValue()
		if currentValue != newValue || newValue == true {
			if newValue == true {
				b.buttonDown <- true
				var d time.Duration = time.Duration(time.Duration(b.repeatInterval-10) * time.Millisecond)
				time.Sleep(d)
			} else {
				b.buttonUp <- false
				currentPressedTime = 0
			}
			currentValue = newValue
		}
		time.Sleep(10 * time.Millisecond)
	}
}
