#!/bin/bash
kubectl delete -f bookinfo-gateway.yaml --context master
kubectl delete -f destination-rule.yaml --context master

# test 1
kubectl delete -f productpage.yaml --context eks-keti-cluster1
kubectl delete -f details.yaml --context eks-keti-cluster1
kubectl delete -f ratings.yaml --context eks-keti-cluster1
kubectl delete -f reviews-v2.yaml --context eks-keti-cluster1

kubectl delete -f productpage.yaml --context  cluster-master1
kubectl delete -f reviews-v1.yaml --context cluster-master1
