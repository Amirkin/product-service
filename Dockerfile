FROM golang:1.13-alpine
EXPOSE 5300
WORKDIR /1
COPY ./product-service product-service
ENTRYPOINT ./product-service