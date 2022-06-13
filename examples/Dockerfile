FROM golang:1.18-alpine

WORKDIR /app

RUN apk add make

RUN apk add --no-cache aws-cli

COPY ./ ./

RUN go mod download

RUN go build main.go -o ./out/repoMan

EXPOSE 3000

CMD out/repoMan
