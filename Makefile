BUILDDIR ?= dist
OSS ?= linux darwin freebsd windows
ARCHS ?= amd64 arm64
VERSION ?= $(shell git describe --tags `git rev-list -1 HEAD`)

build: $(BUILDDIR)/wgrest

clean:
	rm -rf "$(BUILDDIR)"

install: build

# Generate swagger docs
swagger:
	swag init -g cmd/wgrest-server/main.go -o api/docs --ot json

# Run tests
test:
	go test ./... -v

# Run linter
lint:
	golangci-lint run

define wgrest
$(BUILDDIR)/wgrest-$(1)-$(2): export CGO_ENABLED := 0
$(BUILDDIR)/wgrest-$(1)-$(2): export GOOS := $(1)
$(BUILDDIR)/wgrest-$(1)-$(2): export GOARCH := $(2)
$(BUILDDIR)/wgrest-$(1)-$(2):
	go build \
	-ldflags="-s -w -X main.appVersion=$(VERSION)" \
	-trimpath -v -o "$(BUILDDIR)/wgrest-$(1)-$(2)" \
	cmd/wgrest-server/main.go
endef
$(foreach OS,$(OSS),$(foreach ARCH,$(ARCHS),$(eval $(call wgrest,$(OS),$(ARCH)))))

$(BUILDDIR)/wgrest: $(foreach OS,$(OSS),$(foreach ARCH,$(ARCHS),$(BUILDDIR)/wgrest-$(OS)-$(ARCH)))
	@mkdir -vp "$(BUILDDIR)"

.PHONY: clean build install swagger test lint