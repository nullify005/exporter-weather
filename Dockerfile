FROM golang:1.19.1-alpine3.16 AS builder
ARG TARGETARCH
WORKDIR /app
COPY src/go.mod ./go.mod
COPY src/go.sum ./go.sum
RUN go mod download
COPY src/cmd ./cmd
COPY src/internal ./internal
# RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o /exporter-weather ./cmd/exporter-weather
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -ldflags="-s -w" -installsuffix cgo -v -o /exporter-weather ./cmd/exporter-weather

# FROM builder AS test
# RUN go install honnef.co/go/tools/cmd/staticcheck@latest
# RUN staticcheck cmd/exporter-weather/main.go

FROM scratch
COPY --from=builder /exporter-weather /exporter-weather
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/exporter-weather"]