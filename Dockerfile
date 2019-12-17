FROM alpine:latest

RUN addgroup api && \
  adduser -D -G api api

RUN mkdir /app

COPY ./bin/product-api /app/product-api

ENTRYPOINT [ "/app/product-api" ]