FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o worker

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /src/worker .

CMD ["./worker"]
