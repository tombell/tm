VERSION?=dev
COMMIT=$(shell git rev-parse HEAD | cut -c -8)

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"
MODFLAGS=-mod=vendor
TESTFLAGS=-cover

PLATFORMS:=darwin linux windows

dev:
	@echo building dist/tm
	@go build ${MODFLAGS} ${LDFLAGS} -o dist/tm ./cmd/tm

prod: $(PLATFORMS)

$(PLATFORMS):
	@echo building dist/tm-$@-amd64
	@GOOS=$@ GOARCH=amd64 go build ${MODFLAGS} ${LDFLAGS} -o dist/tm-$@-amd64 ./cmd/tm

test:
	@go test ${MODFLAGS} ${TESTFLAGS} ./...

clean:
	@rm -fr dist

.DEFAULT_GOAL := dev
.PHONY: dev prod $(PLATFORMS) test clean
