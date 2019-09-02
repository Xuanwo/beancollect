SHELL := /bin/bash

.PHONY: all check format　vet lint build install uninstall release clean test

VERSION=$(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build beancollect"
	@echo "  install    to install beancollect to /usr/local/bin/beancollect"
	@echo "  uninstall  to uninstall beancollect"
	@echo "  release    to release beancollect"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"

check: format vet lint

format:
	@echo "go fmt"
	go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet -all ./...
	@echo "ok"

lint:
	@echo "golint"
	golint ./...
	@echo "ok"

build: check
	@echo "build beancollect"
	@mkdir -p ./bin
	@go build -tags netgo -o ./bin/beancollect ./cmd/beancollect
	@echo "ok"

install: build
	@echo "install beancollect to GOPATH"
	@cp ./bin/beancollect ${GOPATH}/bin/beancollect
	@echo "ok"

release:
	@echo "release beancollect"
	@rm ./release/*
	@mkdir -p ./release

	@echo "build for linux"
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/beancollect_v${VERSION}_linux_amd64 .
	@tar -C ./bin/linux/ -czf ./release/beancollect_v${VERSION}_linux_amd64.tar.gz beancollect_v${VERSION}_linux_amd64

	@echo "build for macOS"
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/macos/beancollect_v${VERSION}_macos_amd64 .
	@tar -C ./bin/macos/ -czf ./release/beancollect_v${VERSION}_macos_amd64.tar.gz beancollect_v${VERSION}_macos_amd64

	@echo "build for windows"
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows/beancollect_v${VERSION}_windows_amd64.exe .
	@tar -C ./bin/windows/ -czf ./release/beancollect_v${VERSION}_windows_amd64.tar.gz beancollect_v${VERSION}_windows_amd64.exe

	@echo "ok"

clean:
	@rm -rf ./bin
	@rm -rf ./release
	@rm -rf ./coverage

test:
	@echo "run test"
	@go test -v ./...
	@echo "ok"
