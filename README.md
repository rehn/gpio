gpio
=======

Userfriendly package to use gpio-pins in golangfor your raspberry pi

Initial version 
  * Generic Gpio support
  * Simple Lcd-display ( 2 Rows x 16 chars ) my test enviroment

Next version 
  * Button structure
  * Lcd Section 
    Create section for part of Lcd-display for easy writing and update.




Goal.
  A userfriendly package for easy use of gpio-pins on your raspberry pi in golang.



Short Examples
=======

 #Regular Gpio - Blink led 1 time
  g := gpio.NewGpio("out", 4)params out/in, num of gpio 
  g.SetHigh()
  time.Sleep(1 * time.Second)
  g.SetLow()

#Lcd
 #New lcd parameters ( en int, rw int, rs int, D0 int, D1 int, D2 int, D3 int, D4 int, D5 int, D6 int, D7 int, height int, width int )
  
 #skip D0 to D3 for now only using 4 bits
  
  lcd := gpio.NewLcd(8, 0, 7, 0, 0, 0, 0, 25, 24, 23, 18, 2, 16)  
  lcd.WriteString("Hello World")
  
  # section of lcd
  s := lcd.NewSection("clock", 2, 11, 5)//parameters name,line,position,length
  s.WriteString("00:00:00")



  #button
  btn := gpio.NewButton(17, 200)
  go handleKeypress(&btn)
  




  time.Sleep(20 * time.Second)
  lcd.Clear()
  g.Cleanup()
  lcd.Dispose()
}




func handleKeypress(b *gpio.Button) {
  for {
    select {
    case <-b.ButtonDown:
      log.Print("Button Pressed")
    case <-b.ButtonUp:
      log.Print("Button Released")

    }
  }
}
