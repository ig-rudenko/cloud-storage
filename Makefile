.PHONY: build

build:
	go build -v ./backend/cmd/app

.DEFAULT_GOAL := build