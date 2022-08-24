#!/bin/bash
kubectl create -f service.yaml

for i in {1..5}; do kubectl exec -it httpbin -c httpbin -- curl http://svc-hello-v1.default.svc.cluster.local:8080; sleep 0.5; done


for i in {1..5}; do kubectl exec -it httpbin -c httpbin -- curl http://svc-hello-v2.default.svc.cluster.local:8080; sleep 0.5; done
