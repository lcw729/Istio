#!/bin/bash
kubectl create -f service.yaml

for i in {1..5}; do kubectl exec -it httpbin -c httpbin -- curl http://svc-hello.default.svc.cluster.local:8080; sleep 0.5; done

