Simple go web app using

* [Gin](https://gin-gonic.com/)
* [slog](https://pkg.go.dev/log/slog)
* [OpenTelemetry](https://opentelemetry.io/docs/languages/go/)

The app named `ginapp` just exports otel logs to a separate collector agent. All testable via

```
docker compose up
```
