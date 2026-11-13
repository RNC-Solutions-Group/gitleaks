.DEFAULT_GOAL := help

# ---------------------------------------------------------------------------
# Command contract
#
# Non-mutating: build, test, test-race, test-cover, failfast, verify, lint,
#               version, help, profile (builds a binary; writes only profile data)
# Explicit maintenance: format, tidy, generate, clean
#
# Build and test commands never format, generate, tidy, fetch, or inspect remote
# tags. The release version is resolved from local state unless VERSION is
# supplied explicitly.
# ---------------------------------------------------------------------------

GO ?= go
PKG := github.com/zricethezav/gitleaks/v8
LDFLAGS = -ldflags "-X=$(PKG)/version.Version=$(VERSION)"
COVER = --cover --coverprofile=cover.out

# The minimum version accepted by the embedded default configuration. Use it as
# a valid SemVer base for source builds without an exact local tag.
FALLBACK_VERSION := v8.25.0

# GNU make runs recipes through cmd.exe on Windows and a POSIX shell elsewhere.
# Select the correct null device for the local version probes.
ifeq ($(OS),Windows_NT)
NULL_DEVICE := NUL
else
NULL_DEVICE := /dev/null
endif

GIT_PRESENT := $(shell git rev-parse --is-inside-work-tree 2>$(NULL_DEVICE))
GIT_EXACT_TAG := $(shell git describe --tags --exact-match HEAD 2>$(NULL_DEVICE))
GIT_NEAREST_TAG := $(shell git describe --tags --abbrev=0 2>$(NULL_DEVICE))
GIT_SHORT_SHA := $(shell git rev-parse --short=12 HEAD 2>$(NULL_DEVICE))
GIT_DIRTY := $(shell git status --porcelain --untracked-files=no 2>$(NULL_DEVICE))

# Version precedence:
#   1. VERSION from the command line or environment.
#   2. An exact local tag (with "+dirty" when tracked files are modified).
#   3. The nearest local tag plus "+dev.<commit>" build metadata.
#   4. The default-config minimum plus "+dev.source" when Git is unavailable.
#
# No remote fetch is performed. Build metadata keeps every result valid SemVer,
# which the configuration loader requires.
ifeq ($(strip $(VERSION)),)
  ifeq ($(GIT_PRESENT),true)
    ifneq ($(strip $(GIT_EXACT_TAG)),)
      VERSION := $(strip $(GIT_EXACT_TAG))$(if $(strip $(GIT_DIRTY)),+dirty,)
      VERSION_SOURCE := local exact tag
    else ifneq ($(strip $(GIT_NEAREST_TAG)),)
      VERSION := $(strip $(GIT_NEAREST_TAG))+dev.$(strip $(GIT_SHORT_SHA))
      VERSION_SOURCE := local nearest tag
    else
      VERSION := $(FALLBACK_VERSION)+dev.$(strip $(GIT_SHORT_SHA))
      VERSION_SOURCE := local commit without a tag
    endif
  else
    VERSION := $(FALLBACK_VERSION)+dev.source
    VERSION_SOURCE := source tree without Git metadata
  endif
else
  VERSION := $(strip $(VERSION))
  VERSION_SOURCE := explicitly supplied
endif

.PHONY: help build test test-race test-cover failfast verify lint version \
	format tidy generate clean profile config/gitleaks.toml

help:
	@echo Non-mutating commands:
	@echo   make build       compile the scanner (no formatting, generation, or tidy)
	@echo   make test        run the full test suite without the race detector
	@echo   make test-race   run the full test suite with the race detector
	@echo   make test-cover  run race tests with coverage and open the HTML report
	@echo   make failfast    stop at the first test failure
	@echo   make verify      check source validity and generated configuration
	@echo   make lint        run golangci-lint
	@echo   make version     print the resolved build version
	@echo   Explicit maintenance commands (these modify files or cover artifacts):
	@echo   make format      rewrite Go sources with go fmt
	@echo   make tidy        rewrite go.mod/go.sum with go mod tidy
	@echo   make generate    regenerate config/gitleaks.toml with go generate
	@echo   make clean       remove profiles and test artifacts
	@echo   Override VERSION explicitly, for example: make build VERSION=v8.30.0

build:
	$(GO) build -mod=readonly -trimpath $(LDFLAGS)

test:
	$(GO) test -mod=readonly -count=1 ./...

test-race:
	$(GO) test -mod=readonly -count=1 -race ./...

test-cover:
	$(GO) test -mod=readonly -count=1 -v -race $(COVER) ./...
	$(GO) tool cover -html=cover.out

failfast:
	$(GO) test -mod=readonly -count=1 -failfast ./...

verify:
	$(GO) vet -mod=readonly ./...
	$(GO) test -mod=readonly -count=1 ./cmd/generate/config -run '^TestGeneratedConfigMatchesTracked$$'
ifeq ($(OS),Windows_NT)
	@echo Skipping go mod tidy -diff on this CRLF checkout; Linux CI verifies module consistency.
else
	$(GO) mod tidy -diff
endif

lint:
	golangci-lint run

version:
	@echo $(VERSION)
	@echo version source: $(VERSION_SOURCE)

format:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

generate:
	$(GO) generate ./...

ifeq ($(OS),Windows_NT)
clean:
	@powershell -NoProfile -Command "Remove-Item -Recurse -Force profile -ErrorAction SilentlyContinue; Get-ChildItem -Recurse -File -Include '*.got.*','*.out' | Remove-Item -Force"
else
clean:
	rm -rf profile
	find . -type f -name '*.got.*' -delete
	find . -type f -name '*.out' -delete
endif

profile: build
	./scripts/profile.sh './gitleaks' '.'

# Kept for compatibility with the documented generator invocation.
config/gitleaks.toml:
	$(GO) generate ./...
