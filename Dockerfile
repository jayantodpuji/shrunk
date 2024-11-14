FROM golang:1.22.2-alpine

WORKDIR /app

COPY go.mod ./

COPY *.go ./

RUN go mod download

RUN go mod tidy

RUN go build -o shrunk-api

EXPOSE 8000

CMD ./shrunk-api