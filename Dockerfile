# syntax=docker/dockerfile:1

# ---- build stage -------------------------------------------------------------
# The pinned build toolchain (Go 1.27.4). The module's minimum supported Go
# version is declared in go.mod and is deliberately lower than this pin.
# Build on the native BuildKit platform and cross-compile for the target image;
# this keeps linux/arm64 image builds from running the Go compiler under QEMU.
FROM --platform=$BUILDPLATFORM golang:1.27.4-alpine3.24 AS build

ARG TARGETOS
ARG TARGETARCH

# git is only needed in the build stage: it supplies the version stamp below.
RUN apk add --no-cache git

WORKDIR /go/src/github.com/zricethezav/gitleaks

# Populate the module cache in its own layer so source edits do not re-download.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Version stamping: honour an explicit --build-arg VERSION (release and CI).
# Gitleaks validates the current version as SemVer when loading its config, so
# tagless builds must not use a bare Git SHA or "unknown" as the version. The
# default config requires v8.25.0; use that non-release base plus commit metadata
# for builds without a supplied VERSION. .git stays in the context.
ARG VERSION
RUN set -eu; \
    if [ -n "${VERSION:-}" ]; then \
      version="${VERSION}"; \
    elif tag="$(git describe --tags --exact-match HEAD 2>/dev/null)"; then \
      version="${tag}"; \
    elif sha="$(git rev-parse --short=12 HEAD 2>/dev/null)"; then \
      version="v8.25.0+git.${sha}"; \
    else \
      version="v8.25.0+source"; \
    fi; \
    CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH:-amd64}" go build -trimpath \
      -o /out/gitleaks \
      -ldflags "-s -w -X=github.com/zricethezav/gitleaks/v8/version.Version=${version}"

# ---- shared runtime ----------------------------------------------------------
FROM alpine:3.24 AS runtime-base

# Apply the current Alpine 3.24 security patches during image construction, then
# install only the packages needed at scan time:
#   git + openssh-client  scan git repositories, including SSH remotes
#   bash                  retained so existing shell-based invocations keep working
RUN apk upgrade --no-cache && \
    apk add --no-cache bash git openssh-client

# Repositories and bind-mounted workspaces are often owned by another uid; git
# would otherwise refuse to scan them as "dubious ownership". The default image
# retains this established trust policy.
RUN git config --global --add safe.directory '*'

COPY --from=build /out/gitleaks /usr/bin/gitleaks

ENTRYPOINT ["gitleaks"]

# ---- optional non-root runtime ----------------------------------------------
# Build explicitly with:
#   docker build --target runtime-nonroot -t gitleaks:nonroot .
# A plain `docker build` must select the root stage below, not this one.
FROM runtime-base AS runtime-nonroot

RUN addgroup -S -g 10001 gitleaks && \
    adduser -S -D -H -u 10001 -G gitleaks -h /home/gitleaks -s /sbin/nologin gitleaks && \
    install -d -o gitleaks -g gitleaks -m 0700 /home/gitleaks && \
    HOME=/home/gitleaks git config --global --add safe.directory '*' && \
    chown -R gitleaks:gitleaks /home/gitleaks

ENV HOME=/home/gitleaks
USER 10001:10001

# ---- default root runtime ----------------------------------------------------
# This stage is intentionally last: a plain `docker build` must keep producing
# the historical root-based image. Do not move runtime-nonroot below it.
FROM runtime-base AS runtime
