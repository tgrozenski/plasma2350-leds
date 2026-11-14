// Command plasma2350-leds drives an individually addressable (WS2812/NeoPixel)
// LED strip from a Pimoroni Plasma 2350 (RP2350).
//
// Build and flash with the board in BOOTSEL mode:
//
//	tinygo flash -target=pico2 .
package main

import (
	"image/color"
	"machine"

	"github.com/tgrozenski/plasma2350-leds/internal/power"
)

const (
	dataPin = machine.GP15
	numLEDs = 300
)

const (
	supplyAmps = 5.0
	supplyHeadroom = 0.8
	masterBrightness = 0.25
)

var budget = power.Budget{
	SupplyAmps: supplyAmps,
	Headroom:   supplyHeadroom,
}

func main() {
	strip := newStrip(dataPin, numLEDs, budget)

	for {
		christmas(strip)
	}
}

func christmas(s *strip) {
	greenRGB := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	redRGB := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	for i := 0; i < s.len(); i++ {
		if i % 2 == 0 {
			s.set(i, greenRGB)
		} else {
			s.set(i, redRGB)
		}
	}
	s.show()
}
