// Command plasma2350-leds drives a WS2812 LED strip from a Pimoroni
// Plasma 2350 (RP2350).
//
//	make flash
package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

const (
	// dataPin is the Plasma 2350's level-shifted DAT line on the screw
	// terminal. CLK (GP14) is only used by APA102/DotStar strips.
	dataPin = machine.GP15

	// numLEDs must match the strip actually connected.
	numLEDs = 300

	// brightness caps every channel, 0-255. At 300 LEDs this is what keeps
	// the supply honest: full white at 255 would ask for ~18A.
	brightness = 64
)

func main() {
	dataPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ws := ws2812.NewWS2812(dataPin)
	ws.SetBrightness(brightness)

	leds := make([]color.RGBA, numLEDs)

	for {
		christmas(ws, leds)
	}
}

func christmas(ws ws2812.Device, leds []color.RGBA) {
	greenRGB := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	redRGB := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	for i := range leds {
		if i%2 == 0 {
			leds[i] = greenRGB
		} else {
			leds[i] = redRGB
		}
	}
	ws.WriteColors(leds)
}
