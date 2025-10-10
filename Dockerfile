FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and download dependencies
COPY backend/go.mod ./
RUN go mod download

# Copy source code
COPY backend/*.go ./
COPY frontend/ ./frontend/

# Build
RUN go build -o server .

# Runtime image
FROM alpine:latest
WORKDIR /app

# Copy binary and frontend
COPY --from=builder /app/server .
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080

CMD ["./server"]