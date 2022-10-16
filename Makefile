.PHONY: all clean test build

# Variables
GO			= go
GO_RUN		= $(GO) run
GO_BUILD	= $(GO) build
GO_TEST		= $(GO) test
GO_CLEAN	= $(GO) clean
GO_TOOL		= $(GO) tool
GO_VET		= $(GO) vet
GO_FMT		= $(GO) fmt
GO_GENERATE	= $(GO) generate

MAIN_GO		= main.go
BIN_NAME	= bin/gurl
OUT_DIR		= out/
COVER_DIR	= $(OUT_DIR)cover/
COVER_FILE	= $(COVER_DIR)cover.out
COVER_HTML	= $(COVER_DIR)cover.html


all: clean test build

# Clean
clean:
	$(GO_CLEAN)
	@rm -rf $(OUT_DIR)

# Test
test:
	@mkdir -p $(COVER_DIR)
	$(GO_TEST) -cover ./... -coverprofile=$(COVER_FILE)
	$(GO_TOOL) cover -html=$(COVER_FILE) -o $(COVER_HTML)

# Lint
vet:
	$(GO_VET) ./...
fmt:
	$(GO_FMT) ./...

# Build
run:
	$(GO_RUN) $(MAIN_GO)

build:
	$(GO_BUILD) -o $(OUT_DIR)$(BIN_NAME) $(MAIN_GO)