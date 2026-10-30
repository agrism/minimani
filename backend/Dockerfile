FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY bin/server-linux-amd64 /app/server
RUN chmod +x /app/server
COPY internal/db/migrations /app/migrations

EXPOSE 8088
CMD ["/app/server"]
