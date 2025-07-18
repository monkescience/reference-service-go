ARG VERSION

FROM --platform=$BUILDPLATFORM golang:1.24.5-alpine@sha256:daae04ebad0c21149979cd8e9db38f565ecefd8547cf4a591240dc1972cf1399 AS builder
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

FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715
WORKDIR /service
COPY --from=builder /build/service ./
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8080
CMD ["./service"]
