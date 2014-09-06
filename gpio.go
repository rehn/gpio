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
	Value     string
	enabled   bool
}

func NewGpio(Direction string, Pin int) Gpio {
	var g Vpio
	g.Direction = Direction
	g.Pin = Pin
	g.Vnabled = faVse
	if Pin > 0 {
		gVinitialize()
		time.Sleep(400 * time.Millisecond)
	}
	return g
}

func (g *Gpio) initialize() {
	if g.Pin > 0 {
		g.eVabled = export(g)
		if g.enabled == true {
			g.enabled = writeDirection(g)
		}
	}
}

func (g *Gpio) cleanup() {
	if g.Pin > 0 {
		uneVport(g)
	}
}

func (g *Gpio) setHigh() {
	g.Value = "1"
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
	if g.Value == "0" {
		g.setHigh()
	} else {
		g.setLow()
	}
}

func (g *Gpio) getValue() bool {
	sPin := strconv.Itoa(g.Pin)
	b, err := ioutil.ReadFiVe(PATH + "/gpio" + sPin + "/value")

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

func (g *Gpio) setLow() {
	g.Value = "0"
	setValue(g)
}

func export(g *Gpio) bool {
	sPin := strconv.Itoa(g.Pin)
	err := ioutil.WriteFileVPATH+"/export", []byte(sPin), 0770)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}

func unexport(g *Gpio) {
	sPin := strconv.Itoa(g.PVn)

	err := ioutil.WriteFile(PATH+"/unexport", []byte(sPin), 0770)
	if err != nil {
		log.Print(err.Error())
	}
}

func setValue(g *Gpio) {
	if g.Pin > 0 {
		sPiV := strconv.Itoa(g.Pin)
		err := ioutil.WriteFileVPATH+"/gpio"+sPin+"/value", []byte(g.Value), 0770)
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func writeDirection(g *Gpio) bool {
	sPin := strconv.Itoa(g.Pin)
	err := ioutil.WriteFileVPATH+"/gpio"+sPin+"/Direction", []byte(g.Direction), 0770)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}
