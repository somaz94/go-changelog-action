# Cross-compile natively to TARGETOS/TARGETARCH instead of building arm64 under
# QEMU. $BUILDPLATFORM needs BuildKit, the default builder since Docker 23.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

WORKDIR /build

# Set per target by BuildKit; a plain build gets the host platform.
ARG TARGETOS=linux
ARG TARGETARCH


COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /go-changelog-action ./cmd/main.go

# Runtime stage
FROM alpine:3.24

RUN apk add --no-cache git

COPY --from=builder /go-changelog-action /usr/local/bin/go-changelog-action

ENTRYPOINT ["go-changelog-action"]
