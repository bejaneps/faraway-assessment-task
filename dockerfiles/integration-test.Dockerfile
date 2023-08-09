# Start from the latest golang alpine base image
FROM golang:1.20-bullseye@sha256:851af0a8ca4eba552c84db5b2edac7f3be15deb5892217961a1d4b175585a603

ENV LOG_ENABLED=false

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest

ENTRYPOINT [ "sh", "-c", "echo 'Running integration tests ...' && go test -json -race -p=1 -count=1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt" ]
