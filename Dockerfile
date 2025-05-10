FROM golang:1.23.3 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/chadprogress ./cmd/cp/main.go

FROM alpine:3.19
RUN apk add --no-cache bash

COPY --from=builder /app/config /config
COPY --from=builder /app/.env .env
COPY --from=builder /build/chadprogress /chadprogress

ENTRYPOINT ["/chadprogress"]

CMD ["--config_path=/config/dev.yaml"]