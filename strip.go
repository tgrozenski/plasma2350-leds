package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ws2812"

	"github.com/tgrozenski/plasma2350-leds/internal/power"
)

// strip is a small buffered wrapper around the ws2812 driver: pixels are
// staged in memory and pushed to the hardware in one timing-critical burst.
//
// Animations stage colors at full value and ignore brightness entirely.
// Dimming and the current cap are both applied once, in show().
type strip struct {
	dev    ws2812.Device
	pixels []color.RGBA
	budget power.Budget
}

// newStrip configures pin as an output and returns an all-off strip of n LEDs.
//
// On RP2040/RP2350 the driver drives the pin through PIO and the explicit
// Configure is redundant, but NewWS2812 silently falls back to a bit-banged
// driver if the PIO state machine cannot be claimed, and that path does require
// the pin to be an output already.
func newStrip(pin machine.Pin, n int, b power.Budget) *strip {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	return &strip{
		dev:    ws2812.NewWS2812(pin),
		pixels: make([]color.RGBA, n),
		budget: b,
	}
}

// len reports how many LEDs the strip holds.
func (s *strip) len() int { return len(s.pixels) }

// set stages a color for LED i. Out-of-range indexes are ignored.
func (s *strip) set(i int, c color.RGBA) {
	if i < 0 || i >= len(s.pixels) {
		return
	}
	s.pixels[i] = c
}

// fill stages the same color across every LED.
func (s *strip) fill(c color.RGBA) {
	for i := range s.pixels {
		s.pixels[i] = c
	}
}

// clear stages every LED as off and pushes it immediately.
func (s *strip) clear() {
	s.fill(color.RGBA{})
	s.show()
}

// show applies the master brightness and the current budget, then pushes the
// staged pixels out to the strip.
//
// The combined factor is handed to the driver's own SetBrightness rather than
// being multiplied into a second buffer: the driver already scales each channel
// as it streams, so this avoids copying every pixel again on each frame. The
// staged pixels are left untouched, so repeated show() calls do not compound
// the dimming.
func (s *strip) show() {
	f := masterBrightness * s.budget.Scale(s.pixels, masterBrightness)
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}

	s.dev.SetBrightness(uint8(f * 255))
	s.dev.WriteColors(s.pixels)
}
