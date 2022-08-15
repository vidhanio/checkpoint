FROM golang:1.19-rc-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/bin/bot ./cmd/bot

CMD [ "/app/bin/bot", "-students", "/app/data/students.json", "-guilds", "/app/data/guilds.json" ]
