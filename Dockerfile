FROM golang:1.24.3-alpine3.20 AS builder
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build ${GO_BUILD_ARGS} -o service ./cmd/main.go

FROM alpine:3.20
WORKDIR /service
COPY --from=builder /build/service ./
EXPOSE 8080
CMD ["./service"]