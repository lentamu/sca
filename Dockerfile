FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/sca

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

COPY configs/stub.toml ./configs/stub.toml
COPY scripts ./scripts

RUN chmod +x ./scripts/wait-for-it.sh

EXPOSE 8080

CMD ["./app", "-config", "configs/stub.toml"]
