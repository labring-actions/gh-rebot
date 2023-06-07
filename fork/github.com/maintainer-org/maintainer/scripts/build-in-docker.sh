#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE}")/..

CONTAINER_DIR="/go/src/github.com/maintainer-org/maintainer"

cd ${ROOT}

echo "Build maintainer CLI in Docker..."
docker run --rm \
    -v $(pwd):${CONTAINER_DIR} \
    -e GOPATH=/go \
    -w ${CONTAINER_DIR} \
    golang:1.8 \
    sh -c "go build ."
echo "Build successfully."
cd - > /dev/null
