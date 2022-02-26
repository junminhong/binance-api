FROM golang:1.17.6-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
ADD . /app
RUN go mod download
EXPOSE 8080
