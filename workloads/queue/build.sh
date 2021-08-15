#!/bin/bash

BASE_DIR="$(realpath $(dirname "$0"))"
mkdir -p $BASE_DIR/bin

export CGO_ENABLED=0

( cd $BASE_DIR && \
    go build -o bin/main main.go && \
    go build -o bin/init_queues tools/init_queues.go && \
    go build -o bin/benchmark tools/benchmark.go
)
