package gpio

import (
	"log"
	"strings"
	"time"
)

const (
	latchdelay = time.Nanosecond * 50
	HIGH_BITS  = 4
	LOW_BITS   = 0
)

var linestart [4]int = [4]int{128, 192, 148, 212}

type LcdData struct {
	DataMode bool
	Data     []byte
}

type Section struct {
	lcd      *Lcd
	line     int
	position int
	width    int
	value    string
	lcdPos   byte
}

type Lcd struct {
	En       Gpio
	Rs       Gpio
	Rw       Gpio
	D0       Gpio
	D1       Gpio
	D2       Gpio
	D3       Gpio
	D4       Gpio
	D5       Gpio
	D6       Gpio
	D7       Gpio
	line     int
	position int
	height   int
	width    int
	sections map[string]Section
	abort    chan bool
	queue    chan []LcdData
}

func newLcd(en int, rw int, rs int, D0 int, D1 int, D2 int, D3 int, D4 int, D5 int, D6 int, D7 int, height int, width int) Lcd {
	var lcd Lcd
	lcd.En = newGpio("out", en)
	lcd.Rw = newGpio("out", rw)
	lcd.Rs = newGpio("out", rs)
	lcd.D0 = newGpio("out", D0)
	lcd.D1 = newGpio("out", D1)
	lcd.D2 = newGpio("out", D2)
	lcd.D3 = newGpio("out", D3)
	lcd.D4 = newGpio("out", D4)
	lcd.D5 = newGpio("out", D5)
	lcd.D6 = newGpio("out", D6)
	lcd.D7 = newGpio("out", D7)
	lcd.height = height
	lcd.width = width
	lcd.sections = make(map[string]Section)
	lcd.queue = make(chan []LcdData, 200)
	lcd.abort = make(chan bool, 1)
	lcd.initialize()
	go lcd.startLcdWorker()
	return lcd
}

//Queueworker
func (l *Lcd) startLcdWorker() {
	for {
		w := <-l.queue
		for _, k := range w {
			if k.DataMode {
				l.Rs.setHigh()
			}
			for _, c := range k.Data {
				l.writeByte(c)
			}
			if k.DataMode {
				l.Rs.setLow()
			}

		}

	}
}

// set pins and start lcd writer
func (l *Lcd) initialize() {

	l.Rs.setLow()
	l.writeCommandByte8(0x30)
	l.writeCommandByte8(0x30)
	l.writeCommandByte8(0x20)
	l.writeByte(0x08)
	l.writeByte(0x01)
	l.writeByte(0x0C)
	log.Print("initialize Complete")
}

func (l *Lcd) newSection(name string, line int, position int, width int) Section {
	b := byte(linestart[line-1] + position)
	l.sections[name] = Section{lcd: l, line: line, position: position, width: width, lcdPos: b}

	return l.sections[name]
}
func (l *Lcd) getSection(name string) Section {
	return l.sections[name]
}

func (s *Section) writeString(value string) {

	s.value = value
	if len(value) > s.width {
		value = value[:s.width]
	} else if len(value) < s.width {
		value += strings.Repeat(" ", s.width-len(value))
	}
	s.lcd.queue <- []LcdData{LcdData{DataMode: false, Data: []byte{s.lcdPos}}, LcdData{DataMode: true, Data: []byte(value)}}

}

// clear all pins used by Lcd
func (l *Lcd) dispose() {
	l.Rs.cleanup()
	l.Rw.cleanup()
	l.En.cleanup()
	l.D0.cleanup()
	l.D1.cleanup()
	l.D2.cleanup()
	l.D3.cleanup()
	l.D4.cleanup()
	l.D5.cleanup()
	l.D6.cleanup()
	l.D7.cleanup()
}

func (l *Lcd) enable() {
	l.En.setHigh()

	l.En.setLow()
}

// timing method
func (l *Lcd) latch() {
	time.Sleep(latchdelay)
}

// set col
func (l *Lcd) SetPosition(line int, position int) {
	l.line = line
	l.position = position
	b := byte(linestart[l.line-1] + position)

	l.queue <- []LcdData{LcdData{DataMode: false, Data: []byte{b}}}

}

// write string to lcd
func (l *Lcd) writeCommandByte8(ch byte) {
	l.Rs.setLow()
	l.D4.setLow()
	l.D5.setLow()
	l.D6.setLow()
	l.D7.setLow()
	bitArr := byteToBitArray(ch)
	l.writeBits(bitArr, HIGH_BITS)
	l.enable()
}

func (l *Lcd) writeBits(bits [8]uint, row int) {
	var startBit int = row
	dList := []Gpio{l.D4, l.D5, l.D6, l.D7}
	for i := 0; i < 4; i++ {
		b := (int(bits[startBit+i]) > 0)
		dList[i].setValue(b)
	}
}

// write string to lcd
func (l *Lcd) writeByte(ch byte) {
	bitArr := byteToBitArray(ch)
	l.D4.setLow()
	l.D5.setLow()
	l.D6.setLow()
	l.D7.setLow()
	l.writeBits(bitArr, HIGH_BITS)
	l.enable()
	l.writeBits(bitArr, LOW_BITS)
	l.enable()
	//time.Sleep(1 * time.Nanosecond)
}

// write string to lcd
func (l *Lcd) writeString(text string) {

	l.queue <- []LcdData{LcdData{DataMode: true, Data: []byte(text)}}
	// l.Rs.setHigh()
	// for _, c := range text {
	// 	l.writeByte(byte(c))
	// }
	// l.Rs.setLow()
}

// Clear screen
func (l *Lcd) clear() {
	l.queue <- []LcdData{LcdData{DataMode: false, Data: []byte{0x01}}}
}

func byteToBitArray(b byte) [8]uint {
	var a [8]uint
	for i := uint(0); i < 8; i++ {
		v := (b & (1 << i) >> i)
		a[i] = uint(v)
	}
	return a
}
