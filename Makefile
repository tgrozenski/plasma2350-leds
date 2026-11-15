# Build and flash the Plasma 2350 (RP2350) without touching the board.
#
# Why not `tinygo flash`? Its pico2 target uses the "msd" method, which copies
# a UF2 onto the BOOTSEL drive. On this Mac the board enters BOOTSEL correctly
# and macOS binds the whole USB mass-storage stack, but never publishes a BSD
# disk node, so no volume ever mounts and the copy has nothing to target.
#
# So we drive it over PICOBOOT with picotool instead:
#   1. 1200-baud touch on the CDC port -> firmware reboots into BOOTSEL
#   2. picotool load -x                -> writes flash and restarts the app
#
# No BOOT button, no mounted drive.

TARGET     := pico2
BUILD_DIR  := build
UF2        := $(BUILD_DIR)/firmware.uf2
SOURCES    := $(shell find . -name '*.go' -not -path './$(BUILD_DIR)/*') go.mod go.sum

# Re-evaluated on each use, so it still resolves after a reset.
# Override with: make flash PORT=/dev/cu.usbmodemXXXX
PORT       ?= $(firstword $(wildcard /dev/cu.usbmodem*))

# Seconds to wait for the board to re-enumerate in BOOTSEL.
BOOT_WAIT  := 10

.PHONY: all build flash bootsel monitor test clean help

all: build

help: ## show this help
	@grep -hE '^[a-z]+:.*##' $(MAKEFILE_LIST) | sed -e 's/:[^#]*## /|/' | awk -F'|' '{printf "  %-10s %s\n", $$1, $$2}'

build: $(UF2) ## compile the current directory to a UF2

$(UF2): $(SOURCES)
	@mkdir -p $(BUILD_DIR)
	tinygo build -target=$(TARGET) -o $@ .

flash: $(UF2) ## build, reset into BOOTSEL if needed, and load onto the board
	@if picotool info >/dev/null 2>&1; then \
		echo "==> board already in BOOTSEL"; \
	else \
		if [ -z "$(PORT)" ]; then \
			echo "!! no /dev/cu.usbmodem* found and board is not in BOOTSEL."; \
			echo "   Plugged in? Or override: make flash PORT=/dev/cu.usbmodemXXXX"; \
			exit 1; \
		fi; \
		echo "==> resetting $(PORT) into BOOTSEL (1200-baud touch)"; \
		stty -f $(PORT) 1200 >/dev/null 2>&1 || true; \
		for _ in $$(seq 1 $(BOOT_WAIT)); do \
			picotool info >/dev/null 2>&1 && break; \
			sleep 1; \
		done; \
		picotool info >/dev/null 2>&1 || { \
			echo "!! board did not enter BOOTSEL within $(BOOT_WAIT)s"; exit 1; }; \
	fi; \
	echo "==> loading $(UF2)"; \
	picotool load -x $(UF2)

bootsel: ## drop the board into BOOTSEL and stop
	@if picotool info >/dev/null 2>&1; then \
		echo "already in BOOTSEL"; \
	elif [ -n "$(PORT)" ]; then \
		echo "==> resetting $(PORT) into BOOTSEL"; \
		stty -f $(PORT) 1200 >/dev/null 2>&1 || true; \
	else \
		echo "!! no board found"; exit 1; \
	fi

monitor: ## open the board's serial output
	tinygo monitor -target=$(TARGET)

test: ## run host-side tests (no hardware needed)
	go test ./internal/...

clean: ## remove build artifacts
	rm -rf $(BUILD_DIR)
