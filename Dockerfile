FROM golang:1.22-alpine

WORKDIR /app

COPY . /app

RUN go mod tidy

RUN go build -o mnc ./main.go

RUN ls -l .

EXPOSE 8080

ENTRYPOINT ["./mnc"]