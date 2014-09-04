// delay measures the lag between toggling an output pin and getting
// the interrupt that it was toggled, via a second, input pin.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kzyapkov/gpio"
)

var inN = flag.Int("in", 1, "input pin number")
var outN = flag.Int("out", 2, "output pin number")

func main() {
	flag.Parse()
	log.Printf("Using pins %d as input and %d as output\n", *inN, *outN)
	in, err := gpio.OpenPin(*inN, gpio.ModeInput)
	if err != nil {
		log.Panicf("Unable to open pin %d: %v", *inN, err)
	}
	out, err := gpio.OpenPin(*outN, gpio.ModeOutput)
	if err != nil {
		log.Panicf("Unable to open pin %d: %v", *outN, err)
	}
	log.Print("Pins opened...\n")

	signals := make(chan os.Signal, 10)
	signal.Notify(signals, os.Interrupt, os.Kill)

	triggers := make(chan gpio.Edge, 120)
	in.BeginWatch(gpio.EdgeRising, func() {
		go func() { triggers <- gpio.EdgeFalling }()
	})

	var sum, max time.Duration
	var count uint64
	for {
		count++
		log.Print("Start!\n")
		start := time.Now()
		out.Set()
		timeout := time.After(500 * time.Millisecond)
		select {
		case <-signals:
			average := time.Duration(uint64(sum) / count)
			log.Printf("Average delay: %v\n", average)
			log.Printf("Max delay:     %v\n", max)
			return
		case <-triggers:
			took := time.Now().Sub(start)
			out.Clear()
			cleared := time.Now().Sub(start)
			log.Printf("Signal took %v, cleared in %v\n", took, cleared)
			sum += cleared
			if cleared > max {
				max = cleared
			}
		case t := <-timeout:
			took := t.Sub(start)
			log.Printf("Timed out after %v", took)
		}
		out.Clear()
		time.Sleep(500 * time.Millisecond)
	}

}
