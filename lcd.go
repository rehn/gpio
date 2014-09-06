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

func NewLcd(en int, rw int, rs int, D0 int, D1 int, D2 int, D3 int, D4 int, D5 int, D6 int, D7 int, height int, width int) Lcd {
	var lcd Lcd
	lcd.En = NewGpio("out", en)
	lcd.Rw = NewGpio("out", rw)
	lcd.Rs = NewGpio("out", rs)
	lcd.D0 = NewGpio("out", D0)
	lcd.D1 = NewGpio("out", D1)
	lcd.D2 = NewGpio("out", D2)
	lcd.D3 = NewGpio("out", D3)
	lcd.D4 = NewGpio("out", D4)
	lcd.D5 = NewGpio("out", D5)
	lcd.D6 = NewGpio("out", D6)
	lcd.D7 = NewGpio("out", D7)
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
				l.Rs.SetHigh()
			}
			for _, c := range k.Data {
				l.writeByte(c)
			}
			if k.DataMode {
				l.Rs.SetLow()
			}

		}

	}
}

// set pins and start lcd writer
func (l *Lcd) initialize() {

	l.Rs.SetLow()
	l.writeCommandByte8(0x30)
	l.writeCommandByte8(0x30)
	l.writeCommandByte8(0x20)
	l.writeByte(0x08)
	l.writeByte(0x01)
	l.writeByte(0x0C)
	log.Print("initialize Complete")
}

func (l *Lcd) NewSection(name string, line int, position int, width int) Section {
	b := byte(linestart[line-1] + position)
	l.sections[name] = Section{lcd: l, line: line, position: position, width: width, lcdPos: b}

	return l.sections[name]
}
func (l *Lcd) GetSection(name string) Section {
	return l.sections[name]
}

func (s *Section) WriteString(value string) {

	s.value = value
	if len(value) > s.width {
		value = value[:s.width]
	} else if len(value) < s.width {
		value += strings.Repeat(" ", s.width-len(value))
	}
	s.lcd.queue <- []LcdData{LcdData{DataMode: false, Data: []byte{s.lcdPos}}, LcdData{DataMode: true, Data: []byte(value)}}

}

// clear all pins used by Lcd
func (l *Lcd) Dispose() {
	l.Rs.Cleanup()
	l.Rw.Cleanup()
	l.En.Cleanup()
	l.D0.Cleanup()
	l.D1.Cleanup()
	l.D2.Cleanup()
	l.D3.Cleanup()
	l.D4.Cleanup()
	l.D5.Cleanup()
	l.D6.Cleanup()
	l.D7.Cleanup()
}

func (l *Lcd) enable() {
	l.En.SetHigh()

	l.En.SetLow()
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
	l.Rs.SetLow()
	l.D4.SetLow()
	l.D5.SetLow()
	l.D6.SetLow()
	l.D7.SetLow()
	bitArr := byteToBitArray(ch)
	l.writeBits(bitArr, HIGH_BITS)
	l.enable()
}

func (l *Lcd) writeBits(bits [8]uint, row int) {
	var startBit int = row
	dList := []Gpio{l.D4, l.D5, l.D6, l.D7}
	for i := 0; i < 4; i++ {
		b := (int(bits[startBit+i]) > 0)
		dList[i].SetValue(b)
	}
}

// write string to lcd
func (l *Lcd) writeByte(ch byte) {
	bitArr := byteToBitArray(ch)
	l.D4.SetLow()
	l.D5.SetLow()
	l.D6.SetLow()
	l.D7.SetLow()
	l.writeBits(bitArr, HIGH_BITS)
	l.enable()
	l.writeBits(bitArr, LOW_BITS)
	l.enable()
	//time.Sleep(1 * time.Nanosecond)
}

// write string to lcd
func (l *Lcd) WriteString(text string) {

	l.queue <- []LcdData{LcdData{DataMode: true, Data: []byte(text)}}
}

// Clear screen
func (l *Lcd) Clear() {
	l.queue <- []LcdData{LcdData{DataMode: false, Data: []byte{0x01}}}
	time.Sleep(500 * time.Millisecond)
}

func byteToBitArray(b byte) [8]uint {
	var a [8]uint
	for i := uint(0); i < 8; i++ {
		v := (b & (1 << i) >> i)
		a[i] = uint(v)
	}
	return a
}
