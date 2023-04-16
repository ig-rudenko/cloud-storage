.PHONY: build

build:
	go build -v web/backend/cmd/app

.DEFAULT_GOAL := build