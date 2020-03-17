FROM golang:1.13 AS builder

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o app

FROM buildpack-deps:jessie

WORKDIR /usr/local/bin

COPY --from=builder /usr/src/app/app .

EXPOSE 8080

CMD ["./app"]
