FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o vendingApp /app

RUN chmod +x /app/vendingApp

# build a tiny docker image

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/vendingApp /app


CMD ["/app/vendingApp"]

