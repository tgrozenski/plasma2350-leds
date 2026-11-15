package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

// strip is a small buffered wrapper around the ws2812 driver: pixels are
// staged in memory and pushed to the hardware in one timing-critical burst.
type strip struct {
	dev    ws2812.Device
	pixels []color.RGBA
}

// newStrip configures pin as an output and returns an all-off strip of n LEDs.
func newStrip(pin machine.Pin, n int) *strip {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	return &strip{
		dev:    ws2812.NewWS2812(pin),
		pixels: make([]color.RGBA, n),
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

// clear stages every LED as off.
func (s *strip) clear() { s.fill(color.RGBA{}) }

// show pushes the staged pixels out to the strip.
func (s *strip) show() { s.dev.WriteColors(s.pixels) }
