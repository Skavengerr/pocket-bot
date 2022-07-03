FROM golang:1.18 AS builder

COPY . /golang-pocket/
WORKDIR /golang-pocket/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM golang:1.18

WORKDIR /root/

COPY --from=0 /golang-pocket/bin/bot .
COPY --from=0 /golang-pocket/configs configs/

EXPOSE 80

CMD ["./bot"]