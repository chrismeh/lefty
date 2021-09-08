FROM golang:alpine

RUN apk add --no-cache git
WORKDIR /app
ADD . ./

RUN go build -o api ./cmd/web

CMD ["./api"]
