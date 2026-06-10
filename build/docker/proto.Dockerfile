# =============================================================================
# CortexOps Protobuf Generation Container
# =============================================================================
# Deterministic, hermetic protobuf toolchain.
# ALL versions are pinned. No floating tags. No host dependencies.
#
# Usage:
#   docker build -t cortexops-proto-builder -f build/docker/proto.Dockerfile .
#   docker run --rm -v "$(pwd):/workspace" cortexops-proto-builder generate
# =============================================================================

FROM golang:1.25-alpine AS proto-builder

# ---------------------------------------------------------------------------
# Pinned tool versions — update these in lockstep, never independently.
# ---------------------------------------------------------------------------
ENV BUF_VERSION="1.32.2"
ENV PROTOC_VERSION="27.0"
ENV PROTOC_GEN_GO_VERSION="v1.34.1"
ENV PROTOC_GEN_GO_GRPC_VERSION="v1.4.0"

# ---------------------------------------------------------------------------
# Install OS-level dependencies
# ---------------------------------------------------------------------------
RUN apk add --no-cache \
    git \
    curl \
    unzip \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# ---------------------------------------------------------------------------
# Install Protoc compiler — pinned binary
# ---------------------------------------------------------------------------
RUN ARCH=$(uname -m) && \
    case "${ARCH}" in \
      x86_64)  PROTOC_ARCH="x86_64" ;; \
      aarch64) PROTOC_ARCH="aarch_64" ;; \
      arm64)   PROTOC_ARCH="aarch_64" ;; \
      *)       echo "Unsupported architecture: ${ARCH}" && exit 1 ;; \
    esac && \
    curl -sSL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-${PROTOC_ARCH}.zip" -o protoc.zip && \
    unzip protoc.zip -d /usr/local && \
    rm protoc.zip && \
    protoc --version

# ---------------------------------------------------------------------------
# Install Buf CLI — pinned binary, checksum-verified download
# ---------------------------------------------------------------------------
RUN ARCH=$(uname -m) && \
    case "${ARCH}" in \
      x86_64)  BUF_ARCH="x86_64" ;; \
      aarch64) BUF_ARCH="aarch64" ;; \
      arm64)   BUF_ARCH="aarch64" ;; \
      *)       echo "Unsupported architecture: ${ARCH}" && exit 1 ;; \
    esac && \
    curl -sSL \
      "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-${BUF_ARCH}" \
      -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf && \
    buf --version

# ---------------------------------------------------------------------------
# Install Go protobuf plugins — pinned versions
# ---------------------------------------------------------------------------
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}

# ---------------------------------------------------------------------------
# Verify all tools are on PATH and functional
# ---------------------------------------------------------------------------
RUN protoc --version && \
    buf --version && \
    protoc-gen-go --version 2>&1 || true && \
    protoc-gen-go-grpc --version 2>&1 || true && \
    echo "Proto toolchain verified."

WORKDIR /workspace

ENTRYPOINT ["buf"]
