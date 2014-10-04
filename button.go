package gpio

import "time"
import "log"

type Button struct {
	Gp             Gpio
	ButtonDown     chan bool
	ButtonUp       chan bool
	RepeatInterval int64
}

func NewButton(gpioPin int, repeatIterval int64) Button {
	g := NewGpio("in", gpioPin)
	btn := Button{Gp: g, RepeatInterval: repeatIterval, ButtonDown: make(chan bool, 1), ButtonUp: make(chan bool, 1)}
	go btn.buttonWatcher()
	return btn
}

func (b Button) buttonWatcher() {
	currentValue := b.Gp.GetValue()

	for {
		newValue := b.Gp.GetValue()
		if currentValue != newValue || newValue == true {
			if newValue == false {
				b.ButtonUp <- true
			} else {
				b.ButtonDown <- true
				var d time.Duration = time.Duration(time.Duration(b.RepeatInterval-10) * time.Millisecond)
				time.Sleep(d)
			}

			currentValue = newValue
		}
		time.Sleep(10 * time.Millisecond)
	}
}
