FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/importer cmd/importer/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/importer /app/importer
COPY --from=builder /app/internal/db/migrations /app/migrations

EXPOSE 8088
CMD ["/app/server"]

