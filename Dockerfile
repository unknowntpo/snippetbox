FROM golang:1.14.2-alpine3.11

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make

EXPOSE 4000

CMD /app/web
