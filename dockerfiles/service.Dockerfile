FROM golang:1.20-bullseye@sha256:851af0a8ca4eba552c84db5b2edac7f3be15deb5892217961a1d4b175585a603 as builder

ARG SERVICE
ENV SERVICE=${SERVICE}

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
  --mount=type=cache,target=/var/lib/apt,sharing=locked \
  apt update && apt install -y bash-static

WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
  go build -o app ./cmd/${SERVICE}

# final image
FROM golang:1.20-bullseye@sha256:851af0a8ca4eba552c84db5b2edac7f3be15deb5892217961a1d4b175585a603 as base

WORKDIR /

COPY --from=builder /build/app ./app
COPY --from=builder /bin/bash-static /bin/bash

ENTRYPOINT ["./app"]
