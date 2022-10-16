FROM golang:1.19.2-alpine3.16 AS builder
RUN apk --no-cache add build-base
ARG TARGETARCH
WORKDIR /src
COPY src/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -ldflags="-s -w" -installsuffix cgo -v -o /exporter-weather .

FROM builder AS test
RUN go test -v ./...

FROM scratch AS final
COPY --from=builder /exporter-weather /exporter-weather
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 2112/tcp
CMD ["/exporter-weather"]