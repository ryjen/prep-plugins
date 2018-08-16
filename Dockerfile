FROM alpine

ARG CMAKE_DEFINES

RUN apk update

RUN apk add \
    bash \
    tar \
    unzip \
    gzip \
    bzip2 \
    curl \
    g++ \
    make \
    git \
    cmake \
    autoconf \
    automake \
    bats

COPY . /usr/src

RUN mkdir -p /usr/src/docker-build

WORKDIR /usr/src/docker-build

RUN cmake ${CMAKE_DEFINES} ..

RUN make

CMD "make", "test", "ARGS=-V"

