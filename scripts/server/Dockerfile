FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o server cmd/server/main.go

CMD ["/app/server", "-port", "3000"]
