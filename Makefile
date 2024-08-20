APP_NAME := go-padding
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)

build:
	@go build -o $(APP_NAME)$(APP_EXT)

install:
	go install

test:
	@go test

.PHONY: build install test
