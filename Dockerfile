FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY backend/go.mod ./
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