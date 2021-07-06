FROM golang:1.16-alpine AS build
ENV GO111MODULE=on
RUN apk add --no-cache curl git ca-certificates
RUN go get github.com/markbates/pkger/cmd/pkger
WORKDIR /go/src/mario
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN \
  pkger && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mario cmd/mario/main.go

# Note: the two `RUN true` commands appear to be necessary because of
# https://github.com/moby/moby/issues/37965
FROM golang:1.16-alpine
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN true
COPY --from=build /go/src/mario/mario .
RUN true
COPY --from=build /go/src/mario/config ./config
ENTRYPOINT ["./mario"]
CMD ["--help"]
