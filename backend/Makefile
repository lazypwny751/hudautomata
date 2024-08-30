PREFIX := /usr
BUILDD := $(PREFIX)/bin
.DEFAULT_GOAL := build
BINTAR := server.exe

setup:
	mkdir -p $(BUILDD)

build: setup
	go build -o $(BUILDD)/$(BINTAR) src/main.go

clean:
	rm -rf $(BUILDD)/$(BINTAR)

.PHONY: build
