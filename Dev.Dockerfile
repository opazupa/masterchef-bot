FROM golang:1.13

WORKDIR /go/masterchef-bot

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o ./bin ./cmd/app/" -command="./bin/app"