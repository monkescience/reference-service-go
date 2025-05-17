FROM golang:1.24.3-alpine3.20@sha256:9f98e9893fbc798c710f3432baa1e0ac6127799127c3101d2c263c3a954f0abe AS builder
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build ${GO_BUILD_ARGS} -o service ./cmd/main.go

FROM alpine:3.20@sha256:de4fe7064d8f98419ea6b49190df1abbf43450c1702eeb864fe9ced453c1cc5f
WORKDIR /service
COPY --from=builder /build/service ./
EXPOSE 8080
CMD ["./service"]