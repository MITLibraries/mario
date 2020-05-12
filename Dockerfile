FROM golang:1.13-alpine

RUN apk add --no-cache curl git ca-certificates
WORKDIR /go/src/mario
COPY . ./
RUN go build

FROM alpine
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/mario/mario .
ENTRYPOINT ["./mario"]
CMD ["--help"]
