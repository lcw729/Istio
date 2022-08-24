#!/bin/bash
istioctl x uninstall --purge --context $1
kubectl delete namespace istio-system --context $1
