package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"

	"strconv"

	"log"

	"github.com/kzyapkov/gpio"
)

var pinNumbers = flag.String("pins", "", "Comma-separated pin numbers to monitor")

type PinEvent struct {
	Number int
	Value  bool
}

func main() {
	flag.Parse()
	var i, n int
	var err error
	var v string
	pins := make([]int, 100) // should be good for anyone

	for _, v = range strings.Split(*pinNumbers, ",") {
		i, err = strconv.Atoi(v)
		if err == nil {
			pins[n] = i
			n++
		}
	}

	if n == 0 {
		log.Fatalf("No pins specified, got %s", *pinNumbers)
	} else {
		log.Printf("Watching %d pins: %#v", n, pins[:n])
	}

	toggle := make(chan PinEvent, 10)
	for _, i = range pins[:n] {
		var thePin gpio.Pin
		var theNum = i
		thePin, err = gpio.OpenPin(n, gpio.ModeInput)
		if err != nil {
			log.Fatal(err)
		}
		defer thePin.Close()
		thePin.BeginWatch(gpio.EdgeBoth, func() {
			toggle <- PinEvent{theNum, thePin.Get()}
		})
	}

	// clean up on exit
	die := make(chan os.Signal)
	signal.Notify(die, os.Interrupt)
	for {
		select {
		case <-die:
			return
		case e := <-toggle:
			log.Printf("Pin %d is now %t", e.Number, e.Value)
		}
	}
	log.Println("Signal received, returning")
}
