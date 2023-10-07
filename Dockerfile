FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

FROM redis

WORKDIR /data

COPY --from=builder /app/main .
COPY --from=builder /app/app/build /data/app/build

ENTRYPOINT ["sh", "-c", "redis-server --daemonize yes && /data/main"]
