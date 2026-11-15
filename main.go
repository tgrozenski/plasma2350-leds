// Command plasma2350-leds drives a WS2812 LED strip from a Pimoroni
// Plasma 2350 (RP2350).
//
//	make flash
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

// Strip config
const (
	dataPin = machine.GP15
	numLEDs = 300
	brightness = 64
)

func main() {
	dataPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ws := ws2812.NewWS2812(dataPin)
	ws.SetBrightness(brightness)

	strip := make([]color.RGBA, numLEDs)
	//go christmas(&ws, strip)
	setAlternating(strip, color.RGBA{R: 0, G: 255, B: 0, A: 255}, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	go pulse(&ws, strip)

	select {}
}

func setAlternating(strip []color.RGBA, colorTwo, colorOne color.RGBA) {
	for i := range strip {
		if i%2 == 0 {
			strip[i] = colorOne
		} else {
			strip[i] = colorTwo
		}
	}
}

func pulse(ws *ws2812.Device, currentStrip []color.RGBA) {
	current, ceiling := brightness, brightness
	floor := brightness/2

	for {
		for current > floor {
			ws.SetBrightness(uint8(current))
			ws.WriteColors(currentStrip)
			time.Sleep(time.Millisecond*50)
			current--
		}
		for current < ceiling {
			ws.SetBrightness(uint8(current))
			ws.WriteColors(currentStrip)
			time.Sleep(time.Millisecond*50)
			current++
		}
	}
}

func christmas(ws *ws2812.Device, currentStrip []color.RGBA) {
	greenRGB := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	redRGB := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	for {
		setAlternating(currentStrip, redRGB, greenRGB)
		ws.WriteColors(currentStrip)
		time.Sleep(time.Second*1)

		setAlternating(currentStrip, greenRGB, redRGB)
		ws.WriteColors(currentStrip)
		time.Sleep(time.Second*1)
	}
}
