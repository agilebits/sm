GH_USER ?= agilebits
NAME = sm
HARDWARE = $(shell uname -m)
VERSION ?= 0.1.0

build:
	mkdir -p build/linux  && GOOS=linux  go build -a -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -a -o build/darwin/$(NAME)

deps:
	go get -u github.com/progrium/gh-release/...
	go get -u github.com/spf13/viper/...

release: build
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(HARDWARE).tgz -C build/linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(HARDWARE).tgz -C build/darwin $(NAME)
	gh-release create $(GH_USER)/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

clean:
	rm -rf build/*

.PHONY: build release deps clean
