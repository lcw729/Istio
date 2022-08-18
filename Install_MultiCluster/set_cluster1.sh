#!/bin/bash
CTX_CLUSTER1=$1

# Set the default network for cluster1
#kubectl --context="${CTX_CLUSTER1}" get namespace istio-system
#kubectl --context="${CTX_CLUSTER1}" label namespace istio-system topology.istio.io/network=network1

# Configure cluster1 as a primary
istioctl install --context="${CTX_CLUSTER1}" -f cluster1/cluster1.yaml

# Install the east-west gateway in cluster1
./cluster1/gen-eastwest-gateway.sh \
    --mesh mesh1 --cluster cluster1 --network network1 | \
    istioctl --context="${CTX_CLUSTER1}" install -y -f -

# Expose the control plane in cluster1
kubectl apply --context="${CTX_CLUSTER1}" -n istio-system -f cluster1/expose-istiod.yaml

# Expose services in cluster1
#kubectl --context="${CTX_CLUSTER1}" apply -n istio-system -f cluster1/expose-services.yaml
