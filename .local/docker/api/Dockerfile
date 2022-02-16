FROM golang:1.17

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0

RUN go install github.com/cespare/reflex@latest

RUN mkdir /go/src/app-dev
VOLUME /go/src/app-dev

WORKDIR /go/src/app-dev

COPY ./hot-reload.sh /hot-reload.sh
RUN chmod +x /hot-reload.sh

ENTRYPOINT /hot-reload.sh