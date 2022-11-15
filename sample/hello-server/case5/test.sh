#!/bin/bash
kubectl delete virtualservice vs-hello
kubectl create -f virtualservice.yaml
for i in {1..20}; do kubectl exec -it httpbin -c httpbin -- curl http://svc-hello.default.svc.cluster.local:8080; sleep 0.5; done