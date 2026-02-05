FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
# COPY go.sum ./ # Ignore go.sum for simplicity in lab
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o sre-app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sre-app .

EXPOSE 8080
CMD ["./sre-app"]
