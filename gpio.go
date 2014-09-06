package gpio

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

const (
	MOCK = false
)

var PATH string = "/sys/class/gpio" // Raspberry pi

type Gpio struct {
	Direction string
	Pin       int
	value     string
	enabled   bool
}

func NewGpio(direction string, pin int) Gpio {
	var g Gpio
	g.Direction = direction
	g.Pin = pin
	g.enabled = false
	if pin > 0 {
		g.initialize()
		time.Sleep(400 * time.Millisecond)
	}
	return g
}

func (g *Gpio) initialize() {
	if g.Pin > 0 {
		g.enabled = export(g)
		if g.enabled == true {
			g.enabled = writeDirection(g)
		}
	}
}

func (g *Gpio) cleanup() {
	if g.Pin > 0 {
		unexport(g)
	}
}

func (g *Gpio) setHigh() {
	g.value = "1"
	setValue(g)
}

func (g *Gpio) setValue(val bool) {
	if val {
		g.setHigh()
	} else {
		g.setLow()
	}
}
func (g *Gpio) toggleValue() {
	if g.value == "0" {
		g.setHigh()
	} else {
		g.setLow()
	}
}

func (g *Gpio) getValue() bool {
	sPin := strconv.Itoa(g.Pin)
	b, err := ioutil.ReadFile(PATH + "/gpio" + sPin + "/value")

	if err != nil {
		log.Print(err.Error())
		return g.value == "1"
	}
	if b[0] == 48 {
		g.value = "0"
		return false
	} else {
		g.value = "1"
		return true
	}

}

func (g *Gpio) setLow() {
	g.value = "0"
	setValue(g)
}

func export(g *Gpio) bool {
	sPin := strconv.Itoa(g.Pin)
	err := ioutil.WriteFile(PATH+"/export", []byte(sPin), 0770)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func unexport(g *Gpio) {
	sPin := strconv.Itoa(g.Pin)

	err := ioutil.WriteFile(PATH+"/unexport", []byte(sPin), 0770)
	if err != nil {
		log.Print(err.Error())
	}
}

func setValue(g *Gpio) {
	if g.Pin > 0 {
		sPin := strconv.Itoa(g.Pin)
		err := ioutil.WriteFile(PATH+"/gpio"+sPin+"/value", []byte(g.value), 0770)
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func writeDirection(g *Gpio) bool {
	sPin := strconv.Itoa(g.Pin)
	err := ioutil.WriteFile(PATH+"/gpio"+sPin+"/direction", []byte(g.Direction), 0770)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}
