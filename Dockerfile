FROM golang:1.17-alpine as builder
WORKDIR /build/
COPY . .
RUN apk add --no-cache make gcc musl-dev git
RUN make build

FROM alpine:latest
LABEL author=Nathan13888
WORKDIR /app

RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/America/Toronto /etc/localtime
RUN addgroup --gid 1500 bot
RUN adduser \
    --disabled-password \
    --home "/app" \
    --ingroup bot \
    --no-create-home \
    --uid 1500 \
    bot

USER bot

COPY --from=builder --chown=bot /build/bin/wcp /app/bot

#VOLUME [ "/app/student.json", "/app/guilds.json", "app/.env" ]

ENTRYPOINT [ "/app/bot"]

