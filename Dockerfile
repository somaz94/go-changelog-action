# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-changelog-action ./cmd/main.go

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache git

COPY --from=builder /go-changelog-action /usr/local/bin/go-changelog-action

ENTRYPOINT ["go-changelog-action"]
