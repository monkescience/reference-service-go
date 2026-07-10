ARG VERSION

FROM --platform=$BUILDPLATFORM golang:1.26.4@sha256:f96cc555eb8db430159a3aa6797cd5bae561945b7b0fe7d0e284c63a3b291609 AS builder
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
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath ${GO_BUILD_ARGS} -ldflags "-X reference-service-go/internal/build.version=${VERSION}" -o reference-service-go ./cmd/reference-service-go

FROM gcr.io/distroless/static-debian12@sha256:20bc6c0bc4d625a22a8fde3e55f6515709b32055ef8fb9cfbddaa06d1760f838
COPY --from=builder /build/reference-service-go /reference-service-go
EXPOSE 8080
CMD ["/reference-service-go"]
