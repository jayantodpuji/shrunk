FROM golang:1.22.2-alpine

WORKDIR /app

COPY . .

# it's either put the envs here or inject in the run command
ENV DB_HOST=shrunk-postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=rahasia
ENV DB_NAME=shrunk

RUN go build -o shrunk-api

EXPOSE 3002

CMD ./shrunk-api