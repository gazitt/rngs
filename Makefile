NAME = rngs
GITHUB_USERNAME = gazitt
TARGET = /usr/local/bin

VERSION = $(shell git describe --tags --abbrev=0)
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
LDFLAGS = -w -s \
	-X 'main.NAME=$(NAME)' \
	-X 'main.VERSION=$(VERSION)'

RELEASE_DIR = releases
ARTIFACTS_DIR = $(RELEASE_DIR)/artifacts/$(VERSION)

BUILD_TARGETS =\
	build-linux-amd64 \
	build-windows-amd64 \

RELEASE_TARGETS =\
	release-linux-amd64 \
	release-windows-amd64 \

.PHONY: clean build install uninstall $(RELEASE_TARGETS) $(BUILD_TARGETS) $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/$(NAME)$(SUFFIX)

build: $(RELEASE_DIR)/$(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)/$(NAME)$(SUFFIX)

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

$(RELEASE_DIR)/$(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)/$(NAME)$(SUFFIX):
	@go build -ldflags="$(LDFLAGS)" \
	-o $(RELEASE_DIR)/$(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)/$(NAME)$(SUFFIX)

all: $(BUILD_TARGETS)

install:
	@cp -i $(RELEASE_DIR)/$(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)/$(NAME)$(SUFFIX) $(TARGET)

uninstall:
	@rm -i $(TARGET)/$(NAME)$(SUFFIX)

release: $(RELEASE_TARGETS)

release-windows-amd64: build-windows-amd64
	@$(MAKE) release-zip GOOS=windows GOARCH=amd64

release-linux-amd64: build-linux-amd64
	@$(MAKE) release-targz GOOS=linux GOARCH=amd64

$(ARTIFACTS_DIR):
	@mkdir -p $(ARTIFACTS_DIR)

release-targz: $(ARTIFACTS_DIR)
	tar -czf $(ARTIFACTS_DIR)/$(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz -C $(RELEASE_DIR) $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)

release-zip: $(ARTIFACTS_DIR)
	cd $(RELEASE_DIR) && zip -9 $(CURDIR)/$(ARTIFACTS_DIR)/$(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)/*

upload-deps:
	go get -v github.com/tcnksm/ghr

release-upload: upload-deps
	ghr -u $(GITHUB_USERNAME) -t ${GITHUB_TOKEN} --draft --replace $(VERSION) $(ARTIFACTS_DIR)

deps:
	go get -v github.com/gazitt/flago
	go get -v github.com/robertkrimen/otto

test: deps
	go test -v

clean:
	-rm -rf $(RELEASE_DIR)/*/*
	-rm -rf $(ARTIFACTS_DIR)/*
