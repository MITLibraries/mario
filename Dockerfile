FROM golang:1.11

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/mario
COPY . .
RUN dep ensure
RUN go install

ENTRYPOINT ["mario"]
CMD ["--help"]
