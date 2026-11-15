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
	"time"

	"github.com/tgrozenski/plasma2350-leds/internal/power"
)

// ---- Hardware configuration -------------------------------------------------

const (
	// dataPin is the Plasma 2350's level-shifted strip data line (DAT).
	// The board's screw terminal DAT is wired to GP15; CLK (GP14) is only
	// used by clocked strips such as APA102/DotStar.
	dataPin = machine.GP15

	// numLEDs must match the strip actually connected.
	numLEDs = 300
)

// ---- Power configuration ----------------------------------------------------

const (
	// supplyAmps is the rated continuous output of the 5V supply, in amps.
	//
	// !! SET THIS TO YOUR ACTUAL PSU RATING !!
	//
	// It is the single most important value in this file. 300 WS2812s at
	// full white want ~18A; this is what stops the firmware from asking for
	// more than the supply can give. Setting it too high does not blow
	// anything up directly, but it lets the rail sag — see internal/power.
	supplyAmps = 5.0

	// supplyHeadroom is the fraction of the rating actually used. Supplies
	// are happiest well below their label, and the label is often optimistic.
	supplyHeadroom = 0.8

	// masterBrightness scales all output before the power budget is applied.
	// Use it for taste; use supplyAmps for safety.
	masterBrightness = 0.25
)

// budget caps what any animation can draw, regardless of what it asks for.
var budget = power.Budget{
	SupplyAmps: supplyAmps,
	Headroom:   supplyHeadroom,
}

func main() {
	strip := newStrip(dataPin, numLEDs, budget)

	for {
		rainbow(strip, 2*time.Second)
	}
}

// rainbow sweeps a hue gradient along the strip for the given duration.
func rainbow(s *strip, d time.Duration) {
	const step = 10 * time.Millisecond
	frames := int(d / step)

	for f := 0; f < frames; f++ {
		for i := 0; i < s.len(); i++ {
			hue := uint8((i*256/s.len() + f*4) % 256)
			s.set(i, hueToRGB(hue))
		}
		s.show()
		time.Sleep(step)
	}
}

// hueToRGB converts a 0-255 hue into a fully saturated RGB color.
func hueToRGB(h uint8) color.RGBA {
	region := h / 43
	rem := (h - region*43) * 6

	switch region {
	case 0:
		return color.RGBA{R: 255, G: rem, B: 0, A: 255}
	case 1:
		return color.RGBA{R: 255 - rem, G: 255, B: 0, A: 255}
	case 2:
		return color.RGBA{R: 0, G: 255, B: rem, A: 255}
	case 3:
		return color.RGBA{R: 0, G: 255 - rem, B: 255, A: 255}
	case 4:
		return color.RGBA{R: rem, G: 0, B: 255, A: 255}
	default:
		return color.RGBA{R: 255, G: 0, B: 255 - rem, A: 255}
	}
}
