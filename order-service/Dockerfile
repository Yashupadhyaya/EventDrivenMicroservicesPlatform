
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /order-service cmd/main.go

EXPOSE 8081
EXPOSE 8082

CMD [ "/user-management-service" ]
