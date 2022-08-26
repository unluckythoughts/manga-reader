FROM golang:1.19-alpine AS base
RUN apk update && \
    apk --no-cache add ca-certificates curl git tzdata dateutils \
      bash build-base sqlite-dev && \
    rm -rf /var/cache/apk/*

ENV CGO_ENABLED=1
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install gotest.tools/gotestsum@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

FROM base AS build

WORKDIR /go/src/github.com/unluckythoughts/manga-reader
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install -mod=vendor -v ./

WORKDIR /go/bin


FROM build AS debug
ENTRYPOINT ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue"]

WORKDIR /go/src/github.com/unluckythoughts/manga-reader

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install -gcflags="all=-N -l" -mod=vendor -v ./...

WORKDIR /go/bin


FROM build AS test
WORKDIR /go/src/github.com/unluckythoughts/manga-reader
ENTRYPOINT ["gotestsum","--","-mod=vendor","-count=1","-timeout=20m","-tags=integration","-p=1","-v"]
CMD [ "./tests/..." ]


FROM build AS debugTest
WORKDIR /go/src/github.com/unluckythoughts/manga-reader
ENTRYPOINT [ "dlv", "--build-flags='-tags=integration -mod=vendor'", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "test"]


# hadolint ignore=DL3007
FROM alpine:latest AS prod
RUN apk update && \
    apk --no-cache add ca-certificates tzdata curl chromium && \
    rm -rf /var/cache/apk/*

CMD ["./manga-reader"]
WORKDIR /app

COPY --from=build /go/bin/manga-reader /app/manga-reader
