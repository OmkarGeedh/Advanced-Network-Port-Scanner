# Stage 1 — Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -x -o scanner main.go

# Stage 2 — Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scanner .

RUN adduser -D scanneruser
USER scanneruser

ENTRYPOINT ["./scanner"]
CMD ["-h"]
