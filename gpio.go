package gpio

import (
	"io/ioutil"
	"log"
	"os"
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
	Value     string
	Enabled   bool
}

func NewGpio(direction string, pin int) Gpio {
	var g Gpio
	g.Direction = direction
	g.Pin = pin
	g.Enabled = false
	if pin > 0 {
		g.Initialize()
		time.Sleep(400 * time.Millisecond)
	}
	return g
}

func (g *Gpio) Initialize() {
	if g.Pin > 0 {
		g.Enabled = export(g)
		if g.Enabled == true {
			g.Enabled = writeDirection(g)
		}
	}
}

func (g *Gpio) Cleanup() {
	if g.Pin > 0 {
		unexport(g)
	}
}

func (g *Gpio) SetHigh() {
	g.Value = "1"
	setValue(g)
}

func (g *Gpio) SetValue(val bool) {
	if val {
		g.SetHigh()
	} else {
		g.SetLow()
	}
}
func (g *Gpio) ToggleValue() {
	if g.Value == "0" {
		g.SetHigh()
	} else {
		g.SetLow()
	}
}

func (g *Gpio) GetValue() bool {
	sPin := strconv.Itoa(g.Pin)
	b, err := ioutil.ReadFile(PATH + "/gpio" + sPin + "/value")

	if err != nil {
		log.Print(err.Error())
		return g.Value == "1"
	}
	if b[0] == 48 {
		g.Value = "0"
		return false
	} else {
		g.Value = "1"
		return true
	}

}

func (g *Gpio) SetLow() {
	g.Value = "0"
	setValue(g)
}

func export(g *Gpio) bool {
	sPin := strconv.Itoa(g.Pin)
	if _, err := os.Stat(PATH + "/gpio" + sPin); os.IsNotExist(err) {
		err := ioutil.WriteFile(PATH+"/export", []byte(sPin), 0770)
		if err != nil {
			log.Print(err.Error())
			return false
		}
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
		err := ioutil.WriteFile(PATH+"/gpio"+sPin+"/value", []byte(g.Value), 0770)
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
