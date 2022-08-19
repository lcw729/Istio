#!/bin/bash

# Install Istio
istioctl install --set profile=demo -y

# Add a namespace label to instruct istio to automatically inject Envoy sidecar proxies.
kubectl label namespace default istio-injection=enabled

# Deploy the sample application
#kubectl apply -f sample/
