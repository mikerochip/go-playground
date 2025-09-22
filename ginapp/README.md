Simple go web app using

* [Gin](https://gin-gonic.com/)
* [slog](https://pkg.go.dev/log/slog)
* [OpenTelemetry](https://opentelemetry.io/docs/languages/go/)

The app is named `ginapp` and can be tested standalone with

```
go run .
```

Standalone will fail attempts to export otel slogs.

OTel slogs are intended to be sent to a collector agent to reduce log noise in the app. App + Collector are testable via

```
docker compose up
```
