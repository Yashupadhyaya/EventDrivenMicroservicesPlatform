
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /inventory-service cmd/main.go

EXPOSE 8083

CMD [ "/inventory-service" ]
