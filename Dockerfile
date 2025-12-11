FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY taskAPI/ ./taskAPI/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./taskAPI/cmd/api-server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser
USER appuser

WORKDIR /home/appuser

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]