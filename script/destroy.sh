#!/bin/bash -ex

kubectl delete -f deployment/kubernetes/csi-node-rbac.yml
kubectl delete -f deployment/kubernetes/csi-node.yml

kubectl delete -f deployment/kubernetes/csi-attacher-rbac.yml
kubectl delete -f deployment/kubernetes/csi-attacher.yml
