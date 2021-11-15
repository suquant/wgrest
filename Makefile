BUILDDIR ?= dist
OSS ?= linux darwin freebsd windows
ARCHS ?= amd64 arm64

build: $(BUILDDIR)/wgrest

clean:
	rm -rf "$(BUILDDIR)"

install: build

define wgrest
$(BUILDDIR)/wgrest-$(1)-$(2): export appVersion=$$(git describe --tags `git rev-list -1 HEAD`)
$(BUILDDIR)/wgrest-$(1)-$(2): export CGO_ENABLED := 0
$(BUILDDIR)/wgrest-$(1)-$(2): export GOOS := $(1)
$(BUILDDIR)/wgrest-$(1)-$(2): export GOARCH := $(2)
$(BUILDDIR)/wgrest-$(1)-$(2):
	go build \
	-ldflags="-s -w -X main.appVersion=$${appVersion}" \
	-trimpath -v -o "$(BUILDDIR)/wgrest-$(1)-$(2)" \
	cmd/wgrest-server/main.go
endef
$(foreach OS,$(OSS),$(foreach ARCH,$(ARCHS),$(eval $(call wgrest,$(OS),$(ARCH)))))

$(BUILDDIR)/wgrest: $(foreach OS,$(OSS),$(foreach ARCH,$(ARCHS),$(BUILDDIR)/wgrest-$(OS)-$(ARCH)))
	@mkdir -vp "$(BUILDDIR)"

go-echo-server:
	openapi-generator generate -g go-echo-server \
		-i openapi-spec.yaml \
		-o . \
		--git-host github.com \
		--git-user-id suquant \
		--git-repo-id wgrest

typescript-axios-client:
	swagger-codegen generate -l typescript-axios \
		-i openapi-spec.yaml \
		-o clients/typeascript-axios

.PHONY: clean build install