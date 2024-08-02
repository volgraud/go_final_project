FROM golang:1.22.5

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo_app

CMD ["./todo_app"]