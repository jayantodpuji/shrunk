FROM golang:1.22.2-alpine

WORKDIR /app

COPY . .

RUN go build -o shrunk-api

EXPOSE 3002

CMD ./shrunk-api