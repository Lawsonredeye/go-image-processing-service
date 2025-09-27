FROM golang:1.25 AS builder

WORKDIR /go_app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /go_app/cmd/api/main ./...

CMD [ "/go_app/cmd/api/main" ]