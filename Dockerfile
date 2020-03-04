FROM golang:latest
RUN mkdir identity-server
ADD .  /identity-server/
WORKDIR /identity-server
RUN go build -o identity-server .
ENV env="test"
CMD ["./identity-server"]