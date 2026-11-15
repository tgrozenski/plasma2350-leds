package power

import (
	"image/color"
	"math"
	"testing"
)

func fill(n int, c color.RGBA) []color.RGBA {
	px := make([]color.RGBA, n)
	for i := range px {
		px[i] = c
	}
	return px
}

var (
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	off   = color.RGBA{}
)

func closeTo(got, want, tol float32) bool {
	return math.Abs(float64(got-want)) <= float64(tol)
}

// The headline number: 300 WS2812s at full white should land near 18A.
func TestEstimateFullWhite300(t *testing.T) {
	got := EstimateMA(fill(300, white), 1)
	want := float32(300*3*MAPerChannelFull + 300*MAQuiescent) // 18000 + 300

	if !closeTo(got, want, 1) {
		t.Fatalf("300 LEDs full white = %.0fmA, want %.0fmA", got, want)
	}
	if got < 18000 {
		t.Errorf("estimate %.0fmA should be at least 18A", got)
	}
}

// Even an all-off strip draws the controllers' quiescent current.
func TestEstimateQuiescentFloor(t *testing.T) {
	got := EstimateMA(fill(300, off), 1)
	if want := float32(300 * MAQuiescent); !closeTo(got, want, 0.01) {
		t.Fatalf("all-off 300 LEDs = %.1fmA, want %.1fmA", got, want)
	}
}

// Brightness must not scale the quiescent floor away.
func TestBrightnessDoesNotScaleQuiescent(t *testing.T) {
	if got := EstimateMA(fill(300, white), 0); got < 300*MAQuiescent {
		t.Fatalf("at brightness 0 got %.1fmA, want >= quiescent floor", got)
	}
}

// A frame already inside budget is passed through untouched.
func TestScaleWithinBudgetIsUnity(t *testing.T) {
	b := Budget{SupplyAmps: 5, Headroom: 0.8} // 4000mA
	if got := b.Scale(fill(300, off), 1); got != 1 {
		t.Fatalf("Scale = %v, want 1", got)
	}
}

// The core safety property: whatever the frame asks for, the scaled draw
// must not exceed the budget.
func TestScaledDrawNeverExceedsBudget(t *testing.T) {
	cases := []struct {
		name  string
		amps  float32
		n     int
		color color.RGBA
	}{
		{"5A supply, 300 white", 5, 300, white},
		{"10A supply, 300 white", 10, 300, white},
		{"2A supply, 300 white", 2, 300, white},
		{"5A supply, 60 white", 5, 60, white},
		{"5A supply, 300 red", 5, 300, color.RGBA{R: 255, A: 255}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := Budget{SupplyAmps: tc.amps, Headroom: 0.8}
			px := fill(tc.n, tc.color)

			s := b.Scale(px, 1)
			if s < 0 || s > 1 {
				t.Fatalf("scale %v out of range [0,1]", s)
			}

			drawn := EstimateMA(px, s)
			if drawn > b.LimitMA()+1 {
				t.Errorf("scaled draw %.0fmA exceeds budget %.0fmA", drawn, b.LimitMA())
			}
		})
	}
}

// A supply too small to even idle the strip yields 0, not a negative scale.
func TestSupplyTooSmallForQuiescent(t *testing.T) {
	b := Budget{SupplyAmps: 0.1, Headroom: 0.8} // 80mA vs 300mA floor
	if got := b.Scale(fill(300, white), 1); got != 0 {
		t.Fatalf("Scale = %v, want 0 for an undersized supply", got)
	}
}

// An unset supply must fail closed rather than allowing full output.
func TestZeroSupplyFailsClosed(t *testing.T) {
	var b Budget
	if got := b.Scale(fill(300, white), 1); got != 0 {
		t.Fatalf("Scale = %v, want 0 for a zero budget", got)
	}
}
