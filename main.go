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
)

const (
	// dataPin is the Plasma 2350's level-shifted strip data line (DAT).
	// The board's screw terminal DAT is wired to GP15; CLK (GP14) is only
	// used by clocked strips such as APA102/DotStar.
	dataPin = machine.GP15

	// numLEDs must match the strip actually connected.
	numLEDs = 60

	// brightness scales every channel. Keep this low while bench testing:
	// 60 LEDs at full white is ~3.6A, well past USB's budget.
	brightness = 0.15
)

func main() {
	strip := newStrip(dataPin, numLEDs)

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
			s.set(i, scale(hueToRGB(hue), brightness))
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

// scale multiplies each color channel by f, clamped to [0,1].
func scale(c color.RGBA, f float32) color.RGBA {
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	return color.RGBA{
		R: uint8(float32(c.R) * f),
		G: uint8(float32(c.G) * f),
		B: uint8(float32(c.B) * f),
		A: c.A,
	}
}
