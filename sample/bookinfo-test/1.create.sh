#!/bin/bash
kubectl create -f bookinfo-gateway.yaml --context master

# test 1
kubectl create -f productpage.yaml --context eks-keti-cluster1
kubectl create -f details.yaml --context eks-keti-cluster1
kubectl create -f ratings.yaml --context eks-keti-cluster1
kubectl create -f reviews-v2.yaml --context eks-keti-cluster1

kubectl create -f productpage.yaml --context  cluster1-master
kubectl create -f reviews-v1.yaml --context cluster1-master
