# Build stage
FROM golang:1.13-alpine AS builder

WORKDIR $GOPATH/src/masterchef-bot/app/
COPY . .

RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/app

# Publish the executable
FROM alpine:latest

COPY --from=builder /go/app /app

EXPOSE 5555
CMD ["/app"]