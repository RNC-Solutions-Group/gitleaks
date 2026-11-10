# syntax=docker/dockerfile:1

# ---- build stage -------------------------------------------------------------
# The pinned build toolchain (Go 1.27.3). The module's minimum supported Go
# version is declared in go.mod and is deliberately lower than this pin.
FROM golang:1.27.3-alpine3.24 AS build

# git is only needed in the build stage: it supplies the version stamp below.
RUN apk add --no-cache git

WORKDIR /go/src/github.com/zricethezav/gitleaks

# Populate the module cache in its own layer so source edits do not re-download.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Version stamping: honour an explicit --build-arg VERSION (what release
# tooling passes). Otherwise use git describe so commits after the latest tag
# aren't mislabeled as that release. It also produces a commit hash for a
# tagless checkout; "unknown" is the last-resort fallback without git metadata.
# Don't use --dirty: .dockerignore omits tracked files from the build context,
# making even clean source appear modified. .git is deliberately included.
ARG VERSION
RUN VERSION="${VERSION:-$(git -c safe.directory='*' describe --tags --always 2>/dev/null || echo unknown)}" && \
    CGO_ENABLED=0 go build -trimpath \
      -o /out/gitleaks \
      -ldflags "-s -w -X=github.com/zricethezav/gitleaks/v8/version.Version=${VERSION}"

# ---- runtime stage -----------------------------------------------------------
FROM alpine:3.24

# Apply the current Alpine 3.24 security patches during image construction, then
# install only the packages needed at scan time:
#   git + openssh-client  scan git repositories, including SSH remotes
#   bash                  retained so existing shell-based invocations keep working
RUN apk upgrade --no-cache && \
    apk add --no-cache bash git openssh-client

# Repositories and bind-mounted workspaces are often owned by another uid; git
# would otherwise refuse to scan them as "dubious ownership".
RUN git config --global --add safe.directory '*'

COPY --from=build /out/gitleaks /usr/bin/gitleaks

ENTRYPOINT ["gitleaks"]
