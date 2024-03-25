FROM golang:1.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ccvalidator cmd/main.go
CMD ["./ccvalidator"]
