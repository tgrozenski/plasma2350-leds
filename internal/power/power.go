// Package power budgets the current drawn by a WS2812 LED strip.
//
// A 5V supply does not "run out of voltage" — it runs out of *current*. Ask a
// supply for more amps than it can source and its output voltage sags instead,
// which shows up as a dim, red-shifted far end (the blue channel browns out
// first), color drift, a flaky data line, and random MCU resets.
//
// This package has no hardware dependencies so it can be tested on the host.
package power

import "image/color"

const (
	// MAPerChannelFull is the current one color channel of one WS2812 draws
	// at value 255. R+G+B at full white is therefore ~60mA per LED.
	MAPerChannelFull = 20.0

	// MAQuiescent is what each WS2812's controller draws even when its LEDs
	// are off. Small per chip, but 300 of them is a real ~300mA floor.
	MAQuiescent = 1.0
)

// Budget describes the supply the strip is allowed to draw from.
type Budget struct {
	// SupplyAmps is the rated continuous output of the 5V supply.
	SupplyAmps float32

	// Headroom is the fraction of that rating actually used, in (0,1].
	// Supplies are happiest well below their label, and the label is
	// often optimistic.
	Headroom float32
}

// LimitMA is the current ceiling in milliamps.
func (b Budget) LimitMA() float32 {
	return b.SupplyAmps * b.Headroom * 1000
}

// EstimateMA approximates the current the given pixels draw once brightness
// has been applied, including the quiescent floor of the controllers.
func EstimateMA(pixels []color.RGBA, brightness float32) float32 {
	var channels float32
	for _, c := range pixels {
		channels += (float32(c.R) + float32(c.G) + float32(c.B)) / 255
	}

	return channels*brightness*MAPerChannelFull + float32(len(pixels))*MAQuiescent
}

// Scale returns a factor in [0,1] to multiply every channel by so the frame
// fits the budget. It returns 1 when the frame already fits.
//
// The quiescent floor cannot be scaled away: if the strip's idle draw alone
// exceeds the budget, Scale returns 0 and the strip stays dark, because the
// supply is too small for this many LEDs at any brightness.
func (b Budget) Scale(pixels []color.RGBA, brightness float32) float32 {
	limit := b.LimitMA()
	if limit <= 0 {
		return 0
	}

	quiescent := float32(len(pixels)) * MAQuiescent
	usable := limit - quiescent
	if usable <= 0 {
		return 0
	}

	draw := EstimateMA(pixels, brightness) - quiescent
	if draw <= usable {
		return 1
	}

	return usable / draw
}
