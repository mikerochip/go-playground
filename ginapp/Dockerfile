FROM golang:1.24.7-alpine AS build

WORKDIR /src
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o app .

FROM gcr.io/distroless/static:latest AS run

COPY --from=build /src/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
