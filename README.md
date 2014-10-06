gpio
=======

Userfriendly package to use gpio-pins in golang for your raspberry pi

Supports 
  * Generic Gpio support
  * Lcd-display
  * Button


Examples
=======

 Gpio
 
      g := gpio.NewGpio("out", 4) 
      g.SetHigh()
      time.Sleep(1 * time.Second)
      g.SetLow()
      g.CleanUp()



 Lcd
 
      NewLcd(en,rs,rw,D0,D1,D2,D3,D4,D5,D6,D7,height,width)
      
      skip D0 to D3 for now only using 4 bits

      lcd := gpio.NewLcd(8, 0, 7, 0, 0, 0, 0, 25, 24, 23, 18, 2, 16)  
  
      lcd.WriteString("Hello World")
  
      //lcd.NewSection(name string,line int,position int,length int)
      lcdSection := lcd.NewSection("clock", 2, 11, 5)//parameters name,line,position,length
      
      lcdSection.WriteString("00:00")

      lcd.Dispose()

Button
    
    //NewButton(pin int,repeatDelay)
    btn := gpio.NewButton(17, 200)
    go func(b btn){
      for {
        select {
          case state <-b.ButtonDown:
            if state == true {
             log.Print("Button is pressed")
            } else {
             log.Print("Button is released")
            }
        }
      }
    }(btn)

