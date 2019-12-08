FROM golang:alpine
RUN mkdir identity-server
ADD .  /identity-server/
WORKDIR /identity-server
RUN go build -o identity-server .
CMD ["./identity-server"]