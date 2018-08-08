# FROM golang:1.11-rc-alpine
FROM ubuntu:bionic

RUN apt update
RUN apt install -y curl

RUN curl -o go.tar.gz https://dl.google.com/go/go1.11beta2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz

ENV GO111MODULE=on
ENV GOPATH=/src/caches/go

RUN apt install -y build-essential git

RUN gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4
RUN curl -o /usr/local/bin/gosu -SL "https://github.com/tianon/gosu/releases/download/1.4/gosu-$(dpkg --print-architecture)" \
    && curl -o /usr/local/bin/gosu.asc -SL "https://github.com/tianon/gosu/releases/download/1.4/gosu-$(dpkg --print-architecture).asc" \
    && gpg --verify /usr/local/bin/gosu.asc \
    && rm /usr/local/bin/gosu.asc \
    && chmod +x /usr/local/bin/gosu

RUN mkdir /src

WORKDIR /src

COPY entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

