FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o server main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
