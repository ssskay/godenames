FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/*.go ./
COPY frontend/ ./frontend/

RUN go build -o server .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080

CMD ["./server"]