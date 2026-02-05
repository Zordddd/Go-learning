FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-w -s -extldflags '-static'" \
        -o /app/api-server \
        ./taskAPI/cmd/api-server/api-server.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/api-server /app/api-server

EXPOSE 8080
CMD ["/app/api-server"]




































#FROM golang:1.25
#
#RUN apt-get update && apt-get install
#
#WORKDIR /app
#
#COPY go.mod go.sum ./
#
#RUN go mod download
#
#COPY . .
#
#RUN go build -o api-server taskAPI/cmd/api-server/api-server.go
#
#EXPOSE 8080
#CMD ["./api-server"]