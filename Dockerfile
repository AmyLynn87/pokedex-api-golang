FROM golang:1.16-alpine

WORKDIR /app

COPY  . .

RUN go mod download

CMD ["go", "run", "./src/main.go"]