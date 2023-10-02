FROM golang:1.21.1

RUN mkdir /app
ADD . /app/
WORKDIR /app/Server
RUN go build -o main .
CMD ["/app/Server/main"]
