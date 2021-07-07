FROM golang:1.16-alpine

RUN apk add --no-cache curl git ca-certificates
RUN go get github.com/markbates/pkger/cmd/pkger
WORKDIR /go/src/mario
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY mario.go .
COPY pkg pkg
COPY cmd cmd
COPY config config
RUN \
  pkger && \
  go build -o mario cmd/mario/main.go

ENTRYPOINT ["./mario"]
CMD ["--help"]
