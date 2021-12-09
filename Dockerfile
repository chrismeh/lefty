FROM golang:alpine AS builder

RUN apk add --no-cache git
WORKDIR /app
ADD . ./

RUN go build -o api ./cmd/web

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/cmd/web/static static
COPY --from=builder /app/cmd/web/templates templates

CMD ["./api"]

