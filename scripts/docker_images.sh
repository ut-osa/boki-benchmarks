#!/bin/bash

ROOT_DIR=`realpath $(dirname $0)/..`

# Use BuildKit as docker builder
export DOCKER_BUILDKIT=1

function build_boki {
    docker build -t zjia/boki:sosp-ae \
        -f $ROOT_DIR/dockerfiles/Dockerfile.boki \
        $ROOT_DIR/boki
}

function push_boki {
    docker push zjia/boki:sosp-ae
}

function build {
    build_boki
}

function push {
    push_boki
}

case "$1" in
build)
    build
    ;;
push)
    push
    ;;
esac
