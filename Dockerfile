FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

# Download all dependencies
RUN go get -d -v ./...

RUN go build -o api

EXPOSE 8000

CMD ["./api"]