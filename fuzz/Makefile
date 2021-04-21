all: run

FUZZ  := go-fuzz
BUILD := go-fuzz-build

build:
	$(BUILD)

run: build
	$(FUZZ) -bin fuzz-fuzz.zip -procs 1 -v 9

