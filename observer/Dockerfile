FROM golang:1.13-alpine as builder

WORKDIR /build

RUN apk --no-cache add upx

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

ENV CGO_ENABLED=0

RUN go test -cover ./...

RUN go build -ldflags="-s -w" ./cmd/observer && \
	upx -4 observer

# ----------------
FROM alpine:3.11

USER 1000:0
WORKDIR /app
COPY --from=builder ./build/observer ./observer

ENTRYPOINT ["/app/observer"]
