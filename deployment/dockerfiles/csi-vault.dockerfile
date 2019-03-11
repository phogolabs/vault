FROM golang:alpine AS builder

ENV PACKAGE github.com/phogolabs/vault
ENV GO111MODULE off
ENV CGO_ENABLED 0
ENV GOOS linux

RUN mkdir -p /go/src/$PACKAGE
RUN mkdir -p /build

ADD . /go/src/$PACKAGE

WORKDIR /go/src/$PACKAGE

RUN go build -v -o /build/csi-vault $PACKAGE/cmd/csi-vault

FROM alpine:3.5

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /app

COPY --from=builder /build/csi-vault /app/

WORKDIR /app

CMD /app/csi-vault
