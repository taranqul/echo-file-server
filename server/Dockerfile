FROM golang:alpine as builder

WORKDIR /app

COPY server.go .

RUN GOOS=linux GOARCH=amd64 go build -o server server.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD [ "./server" ]