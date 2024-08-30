DEV_DIR := $(CURDIR)
APP_REPOSITORY := github.com/t34-dev
APP_NAME := gopad
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)
export GOPRIVATE=$(APP_REPOSITORY)/*

# includes
include .make/get-started.mk
include .make/tag.mk
include .make/test.mk
include .make/help.mk

build:
	@go build -o .bin/$(APP_NAME)$(APP_EXT) cmd/gopad/*

example: build
	@.bin/$(APP_NAME)${APP_EXT} --files "example" -v

install: build
	@cp .bin/$(APP_NAME)${APP_EXT} $(GOPATH)/bin/$(APP_NAME)$(APP_EXT)

.PHONY: install build example
