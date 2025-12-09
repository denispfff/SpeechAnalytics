FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web_server ./cmd/main.go

EXPOSE 8080

ENV TODO_PORT="8080"

ENTRYPOINT ["./web_server"]
CMD ["--help"]

# docker run -d --rm --name my-server-container -p 8080:8080 my-app-image --port=8080 --dbfile=scheduler.db --password=12345
