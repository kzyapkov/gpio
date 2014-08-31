package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
	"flag"

	"github.com/kzyapkov/gpio"
)

var pin = flag.Int("pin", 0, "GPIO # to toggle")
var T = flag.Int("period", 250, "Period to toggle pin at, ms")

func main() {
	flag.Parse()
	fmt.Printf("Opening pin %d and ticking with %d milliseconds", pin, T)
	pin, err := gpio.OpenPin(*pin, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin %d: %s\n", pin, err)
		return
	}

	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nClearing and unexporting the pin.\n")
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	tt := time.Duration(*T) * time.Millisecond
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
