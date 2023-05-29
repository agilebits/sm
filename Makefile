GH_USER ?= agilebits
NAME = sm
HARDWARE = $(shell uname -m)
VERSION ?= 0.8.1

GOOS ?= darwin

build:
	@$(MAKE) build/linux/$(NAME)-amd64
	@$(MAKE) build/linux/$(NAME)-arm64
	@$(MAKE) build/linux/$(NAME)-armhf
	@$(MAKE) build/darwin/$(NAME)-amd64
	@$(MAKE) build/darwin/$(NAME)-arm64

build/darwin/$(NAME)-amd64:
	mkdir -p build/darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/darwin/$(NAME)-amd64

build/darwin/$(NAME)-arm64:
	mkdir -p build/darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/darwin/$(NAME)-arm64

build/linux/$(NAME)-amd64:
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/linux/$(NAME)-amd64

build/linux/$(NAME)-arm64:
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/linux/$(NAME)-arm64

build/linux/$(NAME)-armhf:
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -a -asmflags=-trimpath=/src -gcflags=-trimpath=/src \
										-ldflags "-s -w -X main.Version=$(VERSION)" \
										-o build/linux/$(NAME)-armhf

deps:
	go get -u github.com/progrium/gh-release/...
	go get -u github.com/spf13/viper/...

create-release-artifacts: build
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_amd64.tgz -C build/linux $(NAME)-amd64
	tar -zcf release/$(NAME)_$(VERSION)_linux_arm64.tgz -C build/linux $(NAME)-arm64
	tar -zcf release/$(NAME)_$(VERSION)_linux_armhf.tgz -C build/linux $(NAME)-armhf
	tar -zcf release/$(NAME)_$(VERSION)_darwin_amd64.tgz -C build/darwin $(NAME)-amd64
	tar -zcf release/$(NAME)_$(VERSION)_darwin_arm64.tgz -C build/darwin $(NAME)-arm64

release: create-release-artifacts
	gh-release create $(GH_USER)/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

clean:
	rm -rf build/*

install:
	install build/$(GOOS)/sm $(GOPATH)/bin/

.PHONY: build release deps clean install
