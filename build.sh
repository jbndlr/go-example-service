#!/bin/zsh

DIR=$0:A:h
VERSION=0.0.1
TARGET=$1

build_dev() {
    docker build $DIR \
        --file dev.Dockerfile \
        --tag jbndlr/go-example-service:$VERSION-dev
}

build_dist() {
    docker build $DIR \
        --force-rm \
        --file dist.Dockerfile \
        --tag jbndlr/go-exapmle-service:$VERSION
}

case "$TARGET" in
    dev)
        build_dev
        exit 0
        ;;
    dist)
        build_dist
        exit 0
        ;;
    all)
        build_dev
        build_dist
        exit 0
        ;;
    *)
        echo "No rule for target: $TARGET"
        exit 1
        ;;
esac
