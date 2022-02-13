FROM golang:1.18-rc-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/bin/bot ./cmd/bot

CMD ["/app/bin/bot" "-dotenv", "/app/.env" "-students", "/app/student.json", "-guilds", "/app/guild.json"]