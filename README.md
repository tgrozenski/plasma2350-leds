# plasma2350-leds

Go (TinyGo) firmware for driving an individually addressable WS2812/NeoPixel
LED strip from a **Pimoroni Plasma 2350** (RP2350).

## Hardware

| | |
|---|---|
| Board | Pimoroni Plasma 2350 (RP2350) |
| MCU | RP2350A |
| Strip data (DAT) | **GP15** — level-shifted to 5V on the screw terminal |
| Strip clock (CLK) | GP14 — *unused* by WS2812; only for APA102/DotStar |
| USB ID | `2E8A:000A` |
| Serial port (macOS) | `/dev/cu.usbmodem*` |

The WS2812 driver uses the RP2350's PIO block for hardware-accurate bit timing,
so the data line is not bit-banged.

## Prerequisites

TinyGo is required — stock `go` cannot target the RP2350.

```sh
brew trust tinygo-org/tools
brew install tinygo
```

## Tests

The current-budget math has no hardware dependency, so it runs on the host with
plain `go` — no TinyGo needed:

```sh
go test ./internal/...
```

## Build

```sh
go mod tidy
tinygo build -target=pico2 -o plasma.uf2 .
```

## Flash

The Plasma 2350 does not auto-reset into the bootloader, so enter BOOTSEL manually:

1. Unplug the board.
2. Hold **BOOT/BOOTSEL** while plugging the USB cable back in.
3. Release once `RPI-RP2` appears in `/Volumes`.

Then:

```sh
tinygo flash -target=pico2 .
```

Or drag `plasma.uf2` onto the mounted `RPI-RP2` volume.

## Serial monitor

```sh
tinygo monitor
```

## Power — read this before plugging in 300 LEDs

A 5V supply does not "run out of voltage." It runs out of **current**. Ask for
more amps than the supply can source and its output voltage sags instead:

- the far end of the strip dims and shifts red (blue browns out first)
- colors drift and the data line gets flaky
- the RP2350 randomly resets
- a cheap supply may hiccup, current-limit, or overheat

### The numbers for 300 WS2812Bs

| Condition | Draw |
|---|---|
| Full white, 100% | 300 x 60mA = **~18A** (90W) |
| Quiescent (all off) | ~0.3A — the controllers draw current regardless |
| 25% brightness, white | ~4.5A |

### How this firmware caps it

Set `supplyAmps` in `main.go` to your PSU's actual rating. Every frame is
measured in `internal/power` and scaled down if it would exceed
`supplyAmps * supplyHeadroom`, so no animation can overdraw the supply no
matter what it asks for.

`masterBrightness` is for taste. `supplyAmps` is for safety. Keep them distinct.

### Wiring 300 LEDs

- **Do not** run strip power through USB, and do not run 18A through the
  board's screw terminal. Feed the strip from the PSU directly.
- **Common ground:** PSU ground and Plasma 2350 ground must be tied together,
  or the data signal has no reference and the strip behaves erratically.
- **Power injection:** 5V over 300 LEDs of thin strip copper drops badly. Feed
  +5V/GND from the PSU at both ends, and ideally every ~100 LEDs, or the far
  end stays brown however big the supply is.

## Configuration

Edit the constants at the top of `main.go`:

- `numLEDs` — **must** match the strip actually connected.
- `masterBrightness` — artistic dimming.
- `supplyAmps` — **your PSU rating**; the hard safety cap.

