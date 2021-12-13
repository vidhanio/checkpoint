FROM golang:1.17.5-alpine

WORKDIR /app
COPY . .

RUN go build
CMD [ "./checkpoint" ]