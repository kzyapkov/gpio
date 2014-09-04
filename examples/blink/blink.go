package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kzyapkov/gpio"
)

var pin = flag.Int("pin", 0, "GPIO # to toggle")
var T = flag.Float64("freq", 5, "Frequency in Hz")

func main() {
	flag.Parse()
	fmt.Printf("Opening pin %d and ticking with %.3f Hz", *pin, *T)
	pin, err := gpio.OpenPin(*pin, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin %d: %s\n", pin, err)
		return
	}

	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for _ = range c {
			fmt.Printf("\nClearing and unexporting the pin.\n")
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	// Calculate half period
	periodMs := (1 / *T) * 1000
	tt := time.Duration(periodMs/2) * time.Millisecond
	fmt.Println("")
	for {
		fmt.Print("\rsetting...")
		pin.Set()
		time.Sleep(tt)
		pin.Clear()
		fmt.Print("\rclearing...")
		time.Sleep(tt)
	}
}
