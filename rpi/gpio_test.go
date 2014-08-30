package rpi

import (
	"github.com/kzyapkov/gpio"
)

// assert that rpi.pin implements gpio.Pin
var _ gpio.Pin = new(pin)
