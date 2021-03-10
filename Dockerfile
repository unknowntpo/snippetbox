# build stage
FROM golang:1.14.2-alpine3.11

WORKDIR /app

RUN apk --update add git make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build && make tls-dev

EXPOSE 4000

CMD ["/app/web"]

