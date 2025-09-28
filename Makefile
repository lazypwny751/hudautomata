PREFIX := build
OUT := $(PREFIX)/hudautomata

all: build

build: backend

backend:
	mkdir -p "$(PREFIX)"
	go build -o "$(OUT)" "src/main.go"

clean:
	rm -rf $(OUT)

.PHONY: all build backend clean
