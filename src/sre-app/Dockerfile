FROM golang:alpine AS builder

WORKDIR /app
# Copy source code immediately so go mod tidy can see imports
COPY go.mod *.go ./

# Generate go.sum and download modules inside the container
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o sre-app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sre-app .

EXPOSE 8080
CMD ["./sre-app"]
