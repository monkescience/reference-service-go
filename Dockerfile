FROM golang:1.24.3-alpine3.21@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS builder
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build ${GO_BUILD_ARGS} -o service ./cmd/main.go

FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715
WORKDIR /service
COPY --from=builder /build/service ./
EXPOSE 8080
CMD ["./service"]