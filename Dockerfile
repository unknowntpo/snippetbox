# build stage
FROM golang:1.14.2-alpine3.11 as builder

WORKDIR /app

RUN apk --update add git make

COPY . .

RUN make

# final stage
FROM golang:1.14.2-alpine3.11

WORKDIR /app
COPY --from=builder /app/web .
COPY --from=builder /app/tls/* ./tls/

CMD ["./web"]

