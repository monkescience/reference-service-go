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
COPY ./migrations ./migrations
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -ldflags "-X reference-service-go/internal/build.Version=${VERSION}" -o reference-service-go ./cmd/reference-service-go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /build/reference-service-go /reference-service-go
EXPOSE 8080
CMD ["/reference-service-go"]
