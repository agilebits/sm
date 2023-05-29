GH_USER ?= agilebits
NAME = sm
HARDWARE = $(shell uname -m)
VERSION ?= 0.7.0

GOOS ?= darwin

build:
	@$(MAKE) build/linux/$(NAME)
	@$(MAKE) build/darwin/$(NAME)

build/darwin/$(NAME):
	mkdir -p build/darwin
	CGO_ENABLED=0 GOOS=darwin go build -a -o build/darwin/$(NAME)

build/linux/$(NAME):
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux go build -a -o build/linux/$(NAME)

deps:
	go get -u github.com/progrium/gh-release/...
	go get -u github.com/spf13/viper/...

create-release-artifacts: build
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(HARDWARE).tgz -C build/linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(HARDWARE).tgz -C build/darwin $(NAME)

release: create-release-artifacts
	gh-release create $(GH_USER)/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

clean:
	rm -rf build/*

install:
	install build/$(GOOS)/sm $(GOPATH)/bin/

.PHONY: build release deps clean install
