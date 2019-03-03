#!/bin/bash -e

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker build -t phogolabs/csi-vault -f ./deployment/dockerfile/csi-vault.dockerfile .
docker push phogolabs/csi-vault

