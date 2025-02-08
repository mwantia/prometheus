ARG GO_VERSION=1.23
ARG GO_IMAGE=mcr.microsoft.com/devcontainers/go
ARG ALPINE_VERSION=bookworm

FROM ${GO_IMAGE}:${GO_VERSION}-${ALPINE_VERSION}

# Install dependencies
RUN apt update \
    && apt install -yqq jq protobuf-compiler \
    && apt clean -yqq

# Update non-root user
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Switch to non-root user
USER $USERNAME
WORKDIR /workspace

# Install protoc Go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest &&  \
    go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/Zxilly/go-size-analyzer/cmd/gsa@latest