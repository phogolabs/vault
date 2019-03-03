FROM circleci/golang:latest AS builder

ENV PACKAGE github.com/phogolabs/vault
ENV GO111MODULE on

RUN mkdir -p /go/src/$PACKAGE
RUN mkdir -p /tmp/build

ADD . /go/src/$PACKAGE

WORKDIR /go/src/$PACKAGE

RUN go build -v -o /tmp/build/csi-vault $PACKAGE/cmd/csi-vault 

FROM alpine:3.5

RUN apk add --update ca-certificates
RUN rm -rf /var/cache/apk/* /tmp/*
RUN update-ca-certificates
RUN mkdir -p /app

COPY --from=builder /tmp/build/csi-vault /app/

WORKDIR /app

ENTRYPOINT ./csi-vault
