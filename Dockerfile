FROM golang:1.13-alpine

EXPOSE 5300

WORKDIR /go/src/
COPY . .

RUN go install

ENTRYPOINT /go/bin/product-service