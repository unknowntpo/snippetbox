# build stage
FROM golang:1.14.2-alpine3.11

WORKDIR /app

RUN apk --update add git make

COPY . .

RUN make

EXPOSE 4000

CMD ["/app/web"]

