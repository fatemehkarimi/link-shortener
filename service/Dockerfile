FROM golang:latest as build

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o link-shortener service/main.go


FROM ubuntu:22.04 as run
WORKDIR /app
COPY --from=build /app/link-shortener .

EXPOSE 8080
CMD ["./link-shortener"]
