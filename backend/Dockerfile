FROM golang:1.19-alpine as builder

WORKDIR /project/go-docker/

COPY backend/go.* .
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY backend/ .
RUN go build -o /project/go-docker/build/myapp .

ENTRYPOINT ["sh", "docker-entrypoint.sh"]
