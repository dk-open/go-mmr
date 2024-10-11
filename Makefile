include builder/Makefile-defaults.mk

all: dep

dep:
	go mod tidy
	go mod vendor

clean:
	go clean
	rm -fr vendor

.PHONY: all dep clean install
