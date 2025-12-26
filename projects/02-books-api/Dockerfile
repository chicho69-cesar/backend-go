FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o books-api .


FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs
RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/books-api .

RUN mkdir -p /app/data /app/logs && \
  chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/app/data/books.db
ENV LOG_PATH=/app/logs/api.log

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/configuration || exit 1

CMD ["./books-api"]
