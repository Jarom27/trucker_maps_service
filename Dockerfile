FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app main.go

ENV TRUCKER_DISTRIBUTOR_PORT=8080
CMD ["app"]
