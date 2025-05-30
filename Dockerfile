FROM golang:1.23.4-alpine

WORKDIR /app

# Install git for downloading Air
RUN apk add --no-cache git

# Install Air
RUN go install github.com/air-verse/air@latest


COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]
