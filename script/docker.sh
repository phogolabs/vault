#!/bin/bash -e

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker build -t phogo/csi-vault -f ./deployment/dockerfile/csi-vault.dockerfile .
docker push phogo/csi-vault

