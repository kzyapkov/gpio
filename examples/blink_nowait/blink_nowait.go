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
var T = flag.Float64("freq", 5, "Frequency in Hz")

func loop(die <-chan int, t time.Duration, pin gpio.Pin) {
	tick := time.Tick(t)
	state := false
	fmt.Println("")
	for {
		select {
		case <-tick:
			if state {
				fmt.Print("\rclear        ")
				pin.Clear()
			} else {
				fmt.Print("\rset          ")
				pin.Set()
			}
			state = !state
		case <-die:
			fmt.Println("\ndying...")
			return
		}
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Opening pin %d and ticking with %.3f Hz", *pin, *T)
	pin, err := gpio.OpenPin(*pin, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin %d: %s\n", pin, err)
		return
	}

	// turn the led off on exit
	signals := make(chan os.Signal, 1)
	dying := make(chan int)

	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		closed := false
		for s := range signals {
			fmt.Printf("\nGot %v, exiting...\n", s)
			if !closed {
				close(dying)
			}
			break
		}
	}()

	// Calculate half period
	periodMs := (1 / *T) * 1000
	tt := time.Duration(periodMs / 2) * time.Millisecond

	// run without blocking
	go loop(dying, tt, pin)
	<-dying
	// equivalently, just call loop in this goroutine
	//loop(dying, tt, pin)

	fmt.Print("Clearing and unexporting the pin.\n")
	pin.Clear()
	pin.Close()
}
