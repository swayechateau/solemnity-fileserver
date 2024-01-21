# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fileserver .

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /app/fileserver .
EXPOSE 8080
USER 1000:1000
ENTRYPOINT ["/fileserver"]
