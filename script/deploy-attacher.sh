#!/bin/bash -ex

# kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/driver-registrar/87d0059110a8b4a90a6d2b5a8702dd7f3f270b80/deploy/kubernetes/rbac.yaml
kubectl create -f deployment/kubernetes/csi-attacher-rbac.yml
kubectl create -f deployment/kubernetes/csi-attacher.yml
