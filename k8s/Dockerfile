FROM golang:bullseye
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install sudo
WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go test ./... -cover
RUN go build -o /usr/bin/spinup-user *.go
WORKDIR /

ENTRYPOINT ["/usr/bin/spinup-user"]
