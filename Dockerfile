FROM golang:1.17.3-alpine

WORKDIR /app
COPY . .

RUN go build
CMD [ "./checkpoint" ]