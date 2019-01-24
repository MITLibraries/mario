FROM golang:1.11-alpine

RUN apk add --no-cache curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/mario
COPY Gopkg.* ./
RUN dep ensure -vendor-only
COPY *.go ./
RUN go build

FROM alpine
COPY --from=0 /go/src/mario/mario .
COPY config config/
ENTRYPOINT ["./mario"]
CMD ["--help"]
