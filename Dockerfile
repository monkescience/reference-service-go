ARG VERSION

FROM --platform=$BUILDPLATFORM golang:1.26.1 AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -o service ./cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /build/service /service
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8080
CMD ["/service"]
