FROM golang:1.16-alpine AS build
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
# Note: the two `RUN true` commands appear to be necessary because of
# https://github.com/moby/moby/issues/37965

FROM golang:1.16-alpine
WORKDIR /go/src/mario
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN true
COPY --from=build /go/src/mario/mario .
RUN true
COPY --from=build /go/src/mario/config ./config
run true
COPY --from=0 /go/src/mario/go.mod ./go.mod
RUN true
ENTRYPOINT ["./mario"]
CMD ["--help"]
