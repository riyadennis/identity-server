FROM golang:latest
COPY go.* /identity-server
WORKDIR /identity-server

RUN go mod download

RUN go build -o identity-server .
CMD ["./identity-server"]
