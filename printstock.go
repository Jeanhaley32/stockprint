package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/d2r2/go-hd44780"
	device "github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	"github.com/piquette/finance-go/quote"
)

var (
	stk           = "GOOG"
	address uint8 = 0x27
	bus           = 1
)

func main() {
	i2c, err := i2c.NewI2C(address, bus)
	if err != nil {
		log.Fatalf("Failed to set i2c: %v", err)
	}

	defer i2c.Close()

	lcd, err := device.NewLcd(i2c, device.LCD_16x2)
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go closeRoutine(c, lcd, i2c)

	for {
		quote := GetQuote(stk)
		lcd.ShowMessage(quote[0], device.SHOW_LINE_1)
		lcd.ShowMessage(quote[1], device.SHOW_LINE_2)
		time.Sleep(1 * time.Minute)
	}

}

func closeRoutine(c chan os.Signal, lcd *hd44780.Lcd, i2c *i2c.I2C) {
	for sig := range c {
		log.Printf("captured %v, exiting..", sig)
		lcd.BacklightOff()
		lcd.Clear()
		i2c.Close()
		os.Exit(1)
	}
}

func GetQuote(stock string) [2]string {
	q, err := quote.Get(stock)
	if err != nil {
		log.Fatal(err)
	}
	return [2]string{q.ShortName, fmt.Sprintf("%.2f", q.RegularMarketPrice)}
}
