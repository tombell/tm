NAME = tm

VERSION ?= dev
COMMIT = $(shell git rev-parse HEAD | cut -c -8)

LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"

PLATFORMS := darwin-amd64 darwin-arm64 linux-amd64 linux-arm64

dev:
	@echo "building bin/${NAME}"
	@go build ${LDFLAGS} -o bin/${NAME} ./cmd/${NAME}

prod: $(PLATFORMS)

$(PLATFORMS):
	@echo "building ${NAME}-$@"
	@GOOS=$(word 1,$(subst -, ,$@)) GOARCH=$(word 2,$(subst -, ,$@)) \
		go build ${LDFLAGS} -o bin/${NAME}-$@ ./cmd/${NAME}

watch:
	@while sleep 1; do \
		trap "exit" INT TERM; \
		rg --files --glob '*.go' | \
		entr -c -d -r make dev; \
	done

clean:
	@rm -fr bin

.DEFAULT_GOAL := dev
.PHONY: dev prod $(PLATFORMS) watch clean
