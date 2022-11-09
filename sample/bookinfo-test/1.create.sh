#!/bin/bash

kubectl label namespace default istio-injection=enabled --context master
kubectl label namespace default istio-injection=enabled --context cluster1-master
kubectl label namespace default istio-injection=enabled --context eks-keti-cluster1

kubectl create -f bookinfo-gateway.yaml --context master
kubectl create -f destination-rule.yaml --context master
kubectl create -f bookinfo-gateway.yaml --context eks-keti-cluster1
kubectl create -f destination-rule.yaml --context eks-keti-cluster1
# test 1
kubectl create -f productpage.yaml --context eks-keti-cluster1
kubectl create -f details.yaml --context eks-keti-cluster1
kubectl create -f ratings.yaml --context eks-keti-cluster1
kubectl create -f reviews-v2.yaml --context eks-keti-cluster1

kubectl create -f productpage.yaml --context  cluster1-master
kubectl create -f reviews-v1.yaml --context cluster1-master
