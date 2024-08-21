DEV_DIR := $(CURDIR)
APP_REPOSITORY := github.com/t34-dev
APP_NAME := gopad
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)
export GOPRIVATE=$(APP_REPOSITORY)/*

# includes
include .make/get-started.mk
include .make/tag.mk
include .make/test.mk

build:
	@go build -o $(APP_NAME)$(APP_EXT)

install:
	go build -o $(GOPATH)/bin/$(APP_NAME)$(APP_EXT)

test:
	@go test

.PHONY: build install test
