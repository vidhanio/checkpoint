FROM golang:1.17.4-alpine

WORKDIR /app
COPY . .

RUN go build
CMD [ "./checkpoint" ]