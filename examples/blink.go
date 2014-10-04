package main

import (
	"time"

	"github.com/rehn/gpio"
)

func main() {

	myLed := gpio.NewGpio(4)

	myLed.SetHigh()

	for i := 0; i < 100; i++ {
		myLed.SetValue(!myLed.GetValue())
		time.Sleep(200 * time.Millisecond)
	}

	myLed.SetLow()
	myLed.CleanUp()
}
