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

## Configuration

Edit the constants at the top of `main.go`:

- `numLEDs` — **must** match the strip actually connected.
- `brightness` — scales all channels.

> **Power budget:** USB supplies roughly 0.5A. 60 WS2812s at full white draw
> ~3.6A. Keep `brightness` low unless the strip has its own 5V supply wired to
> the screw terminal.
