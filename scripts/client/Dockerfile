FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o client cmd/client/main.go

CMD ["/app/client", "-host", "localhost", "-port", "3000"]
