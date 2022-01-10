FROM golang:1.17.6-alpine

WORKDIR /app
COPY . .

RUN go build
CMD [ "./checkpoint" ]