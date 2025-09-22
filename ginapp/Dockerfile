FROM golang:1.24.7-alpine AS build

WORKDIR /src
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

# see https://seankhliao.com/blog/12021-01-23-docker-buildx-caching/#-statefull--build-workers
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -o app .

FROM gcr.io/distroless/static:latest AS run

COPY --from=build /src/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
