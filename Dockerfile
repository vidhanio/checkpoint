FROM golang:1.17.2-alpine

WORKDIR /app
COPY . .

RUN go build
CMD [ "./checkpoint" ]