FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN ls -lah

RUN go build -o /docker-gs-ping

RUN ls -lah

EXPOSE 8081

CMD [ "/docker-gs-ping" ]