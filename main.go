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

	go christmas(&ws)
	go pulse(&ws)

	select {}
}

func alternate(strip []color.RGBA, colorTwo, colorOne color.RGBA) {
	for i := range strip {
		if i%2 == 0 {
			strip[i] = colorOne
		} else {
			strip[i] = colorTwo
		}
	}
}

func pulse(ws *ws2812.Device) {
	for {
		time.Sleep(time.Second * 2)
		ws.SetBrightness(16)
		time.Sleep(time.Second * 2)
		ws.SetBrightness(64)
	}

	// current, ceiling := brightness, brightness
	// floor := brightness/4
	//
	// for {
	// 	for current > floor {
	// 		ws.SetBrightness(uint8(current))
	// 		time.Sleep(time.Millisecond*250)
	// 		current--
	// 	}
	// 	for current < ceiling {
	// 		ws.SetBrightness(uint8(current))
	// 		time.Sleep(time.Millisecond*250)
	// 		current++
	// 	}
	// }
}

func christmas(ws *ws2812.Device) {
	greenRGB := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	redRGB := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	ledsOne := make([]color.RGBA, numLEDs)
	alternate(ledsOne, greenRGB, redRGB)

	ledsTwo := make([]color.RGBA, numLEDs)
	alternate(ledsTwo, redRGB, greenRGB)

	for {
		time.Sleep(time.Second*1)
		ws.WriteColors(ledsOne)
		time.Sleep(time.Second*1)
		ws.WriteColors(ledsTwo)
	}
}
