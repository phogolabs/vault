#!/bin/bash -ex

kubectl create -f deployment/kubernetes/csi-node-rbac.yml
kubectl create -f deployment/kubernetes/csi-node.yml

