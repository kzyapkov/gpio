package test

import "github.com/kzyapkov/gpio"

type PinMock struct {
	TheMode    gpio.Mode
	TheError   error
	TheHandler gpio.IRQEvent
	OnEdge     gpio.Edge
	TheValue   bool
}

// gets the current pin mode
func (p *PinMock) Mode() gpio.Mode {
	return p.TheMode
}

// set the current pin mode
func (p *PinMock) SetMode(m gpio.Mode) {
	p.TheMode = m
}

// sets the pin state high
func (p *PinMock) Set() {
	p.TheValue = true
}

// sets the pin state low
func (p *PinMock) Clear() {
	p.TheValue = false
}

// if applicable, closes the pin
func (p *PinMock) Close() error {
	return p.TheError
}

// returns the current pin state
func (p *PinMock) Get() bool {
	return p.TheValue
}

// calls the function argument when an edge trigger event occurs
func (p *PinMock) BeginWatch(edge gpio.Edge, handler gpio.IRQEvent) error {
	p.TheHandler = handler
	p.OnEdge = edge
	return p.TheError
}

// stops watching the pin
func (p *PinMock) EndWatch() error {
	p.TheHandler = nil
	p.OnEdge = gpio.EdgeNone
	return nil
}

// wait for pin state to match boolean argument
func (p *PinMock) Wait(b bool) {
	// hmmm ....
}

// returns the last error state
func (p *PinMock) Err() error {
	return p.TheError
}
