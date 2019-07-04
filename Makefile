SHELL := /bin/bash

.PHONY: all check format　vet lint build install uninstall release clean test coverage

VERSION=$(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_TO_CHECK=$(shell go list ./... | grep -v "/vendor/")
PKGS_TO_CHECK=$(shell go list ./... | grep -vE "/vendor/|/tests/")
INGR_TEST=$(shell go list ./... | grep "/tests/" | grep -v "/utils")

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build beancollect"
	@echo "  install    to install beancollect to /usr/local/bin/beancollect"
	@echo "  uninstall  to uninstall beancollect"
	@echo "  release    to release beancollect"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"
	@echo "  coverage   to test with coverage"

check: format vet lint

format:
	@echo "go fmt, skipping vendor packages"
	@for pkg in ${PKGS_TO_CHECK}; do go fmt $${pkg}; done;
	@echo "ok"

vet:
	@echo "go vet, skipping vendor packages"
	@go vet -all ${DIRS_TO_CHECK}
	@echo "ok"

lint:
	@echo "golint, skipping vendor packages"
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}"); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

build: check
	@echo "build beancollect"
	@mkdir -p ./bin
	@go build -tags netgo -o ./bin/beancollect .
	@echo "ok"

install: build
	@echo "install beancollect to GOPATH"
	@cp ./bin/beancollect ${GOPATH}/bin/beancollect
	@echo "ok"

uninstall:
	@echo "delete /usr/local/bin/beancollect"
	@rm -f /usr/local/bin/beancollect
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
	@go test -v ${PKGS_TO_CHECK}
	@echo "ok"

coverage:
	@echo "run test with coverage"
	@for pkg in ${PKGS_TO_CHECK}; do \
		output="coverage$${pkg#github.com/Xuanwo/beancollect}"; \
		mkdir -p $${output}; \
		go test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			go tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"
