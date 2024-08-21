FROM golang:latest

LABEL maintainer: Reynald Marquez-Gragasin
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
EXPOSE 3306
ENV PORT 5100
RUN go build
CMD ["./goapi"]