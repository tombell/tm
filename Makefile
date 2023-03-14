VERSION?=dev
COMMIT=$(shell git rev-parse HEAD | cut -c -8)

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"
MODFLAGS=-mod=vendor
TESTFLAGS=-cover

PLATFORMS:=darwin linux

dev:
	@echo building dist/tm
	@go build ${MODFLAGS} ${LDFLAGS} -o dist/tm ./cmd/tm

prod: $(PLATFORMS)

$(PLATFORMS):
	@echo building dist/tm-$@-amd64
	@GOOS=$@ GOARCH=amd64 go build ${MODFLAGS} ${LDFLAGS} -o dist/tm-$@-amd64 ./cmd/tm

watch:
	@while sleep 0.1; do \
		trap "exit" SIGINT; \
		find . -type d \( -name vendor \) -prune -false -o -type f \( -name "*.go" \) | entr -d -r make; \
	done

test:
	@go test ${MODFLAGS} ${TESTFLAGS} ./...

clean:
	@rm -fr dist

.DEFAULT_GOAL := dev
.PHONY: dev prod $(PLATFORMS) watch test clean
