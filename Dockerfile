FROM golang:1.19.4-alpine3.16
WORKDIR /var/lib/app

COPY go.mod go.mod
COPY go.sum go.sum
CMD go mod download

COPY ./peripherals ./peripherals
COPY main.go main.go

RUN go build -o app

CMD ./app
