# syntax=docker/dockerfile:1
# DOOM + DOOM II backend stubs (Switch 01008CB01E52E000). Plain HTTP on 23630;
# TLS for the clients' https hosts is terminated in front.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOWORK=off go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/doom-plus-doom-iid ./cmd/doom-plus-doom-iid

FROM gcr.io/distroless/static
COPY --from=build /out/doom-plus-doom-iid /doom-plus-doom-iid
USER nonroot:nonroot
EXPOSE 23630 23631
ENTRYPOINT ["/doom-plus-doom-iid"]
